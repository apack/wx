package wx

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"sync"

	v8 "github.com/apack/wx/internal/js/v8"
	"github.com/apack/wx/internal/ssr"
)

type ViewDOM struct {
	FileName string
	Script   []byte
}

type ViewCSS struct {
	FileName string
	Style    []byte
}

type ViewSSR struct {
	Script []byte
}

type View interface {
	Name() string
	DOM() *ViewDOM
	CSS() *ViewCSS
	SSR() *ViewSSR
}

type Renderer struct {
	view View
	vm   *v8.VM
	ssr  *ssr.SSR
}

type Props map[string]interface{}

func NewRenderer(view View) (*Renderer, error) {
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	ssr, err := ssr.Load(vm, string(view.SSR().Script))
	if err != nil {
		return nil, err
	}
	return &Renderer{
		view: view,
		vm:   vm,
		ssr:  ssr,
	}, nil
}

func (r *Renderer) Close() error {
	r.vm.Close()
	return nil
}

type ssrResult struct {
	HEAD string
	HTML string
}

func (r *Renderer) Render(props Props) ([]byte, error) {
	var dst ssrResult
	err := r.ssr.Execute(&dst, props)
	if err != nil {
		return nil, err
	}
	propsJSON, err := json.Marshal(props)
	if err != nil {
		return nil, err
	}
	html := `<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link rel="stylesheet" href="/styles/` + r.view.CSS().FileName + `">
		` + dst.HEAD + `
	</head>
	<body>
		<div id="app">
		` + dst.HTML + `
		</div>
		<script src="/scripts/` + r.view.DOM().FileName + `"></script>
		<script>
			var view = new __dom__.default({
				target: document.getElementById('app'),
				hydrate: true,
				props: ` + string(propsJSON) + `,
			});
		</script>
	</body>
	</html>`
	return []byte(html), nil
}

type ViewRegistry interface {
	RegisterViews(views ...View) error
}
type closable interface {
	Close() error
}

func makePool[T io.Closer](size int, createFn func() (T, error)) (*pool[T], error) {
	p := &pool[T]{
		slots: make(chan T, size),
	}
	for i := 0; i < size; i++ {
		v, err := createFn()
		if err != nil {
			return nil, err
		}
		p.slots <- v
	}
	return p, nil
}

type pool[T io.Closer] struct {
	slots chan T
}

func (p *pool[T]) get() T {
	return <-p.slots

}
func (p *pool[T]) put(v T) {
	p.slots <- v
}

func (p *pool[T]) Close() error {
	var errs = make([]error, 0)
	for len(p.slots) > 0 {
		v := <-p.slots
		err := v.Close()
		if err != nil {
			errs = append(errs, err)
		}
	}
	if l := len(errs); l > 0 {
		var strArr []string
		for _, err := range errs {
			strArr = append(strArr, err.Error())
		}
		return fmt.Errorf("wx: failed to close %d pool workers:: [%v]", l, strings.Join(strArr, ";"))
	}
	close(p.slots)
	return nil
}

type RenderWorker struct {
	pool *pool[*Renderer]
}

func NewRenderWorker(view View) (*RenderWorker, error) {
	pool, err := makePool(4, func() (*Renderer, error) {
		r, err := NewRenderer(view)
		if err != nil {
			return nil, err
		}
		return r, nil
	})
	if err != nil {
		return nil, err
	}
	return &RenderWorker{
		pool: pool,
	}, nil
}

func (w *RenderWorker) Render(props Props) ([]byte, error) {
	r := w.pool.get()
	defer w.pool.put(r)
	html, err := r.Render(props)
	return html, err
}

func (w *RenderWorker) Close() error {
	return w.pool.Close()
}

type RenderEngine struct {
	wm  map[View]*RenderWorker
	mtx sync.Mutex
}

func (e *RenderEngine) Close() {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	for _, w := range e.wm {
		w.Close()
	}
}

func NewRenderEngine() *RenderEngine {
	return &RenderEngine{
		wm: make(map[View]*RenderWorker),
	}
}

func (e *RenderEngine) RegisterViews(views ...View) error {
	e.mtx.Lock()
	defer e.mtx.Unlock()
	var err error
	for _, view := range views {
		if view == nil {
			return fmt.Errorf("wx: view is nil")
		}
		if _, ok := e.wm[view]; ok {
			return fmt.Errorf("wx: view %s already registered", view.Name())
		}
		e.wm[view], err = NewRenderWorker(view)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *RenderEngine) Render(view View, props Props) ([]byte, error) {
	e.mtx.Lock()
	w, ok := e.wm[view]
	e.mtx.Unlock()
	if !ok {
		return nil, fmt.Errorf("wx: view %q is not registered", view.Name())
	}
	return w.Render(props)
}
