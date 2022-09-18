package wx

import (
	"io/fs"
	"net/http"
)

type Handler interface {
	HandleWeb(ctx *Context) error
}

type HandlerFunc func(ctx *Context) error

func (f HandlerFunc) HandleWeb(ctx *Context) error {
	return f(ctx)
}

type Middleware func(Handler) Handler

type Router struct {
	r   *http.ServeMux
	eng *RenderEngine
	mws []Middleware
}

func NewRouter() *Router {

	return &Router{
		r:   http.NewServeMux(),
		eng: NewRenderEngine(),
	}
}

func (r *Router) Static(files fs.FS) error {
	fs := http.FileServer(http.FS(files))
	r.r.Handle("/static/", fs)
	return nil
}

func (r *Router) RegisterViews(views ...View) error {
	for _, view := range views {
		r.r.HandleFunc("/styles/"+view.CSS().FileName, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "text/css")
			w.Write(view.CSS().Style)
		})
		r.r.HandleFunc("/scripts/"+view.DOM().FileName, func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Content-Type", "application/javascript")
			w.Write(view.DOM().Script)
		})
	}
	return r.eng.RegisterViews(views...)
}

func (r *Router) Use(mws ...Middleware) {
	r.mws = append(r.mws, mws...)
}

func (r *Router) Handle(path string, h Handler) {
	for _, mw := range r.mws {
		h = mw(h)
	}
	r.r.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		ctx := &Context{
			w:       w,
			r:       req,
			eng:     r.eng,
			Context: req.Context(),
		}
		err := h.HandleWeb(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (r *Router) HandleMethod(method, path string, h Handler) {
	for _, mw := range r.mws {
		h = mw(h)
	}
	r.r.HandleFunc(path, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		ctx := &Context{
			w:       w,
			r:       req,
			eng:     r.eng,
			Context: req.Context(),
		}
		err := h.HandleWeb(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (r *Router) HandleFunc(path string, h HandlerFunc) {
	r.Handle(path, h)
}

func (r *Router) HandleMethodFunc(method, path string, h HandlerFunc) {
	r.HandleMethod(method, path, h)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.r.ServeHTTP(w, req)
}
