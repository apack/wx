package compiler

import (
	"crypto/sha1"
	"encoding/base32"
	"fmt"
	"os"
	"path"
	"strings"

	v8 "github.com/apack/wx/internal/js/v8"
	"github.com/apack/wx/internal/ssr"
	"github.com/apack/wx/internal/svelte"
	"github.com/apack/wx/internal/transform"
	esbuild "github.com/evanw/esbuild/pkg/api"
)

type Compiler struct {
	wd  string
	dir string
	vm  *v8.VM
	tm  *transform.Map
}

func NewCompiler(dir string) (*Compiler, error) {
	vm, err := v8.Load()
	if err != nil {
		return nil, err
	}
	compiler, err := svelte.Load(vm)
	if err != nil {
		return nil, err
	}
	tm, err := transform.Load(svelte.NewTransformable(compiler))
	if err != nil {
		return nil, err
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return &Compiler{
		dir: dir,
		wd:  wd,
		vm:  vm,
		tm:  tm,
	}, nil
}

type CompiledView struct {
	Name    string
	CSS     []byte
	CSSHash string
	DOM     []byte
	DOMHash string
	SSR     []byte
	SSRHash string
}

func (c *Compiler) Compile() ([]*CompiledView, error) {
	files, err := os.ReadDir(path.Join(c.dir, "views"))
	if err != nil {
		return nil, err
	}
	encoding := base32.StdEncoding.WithPadding(base32.NoPadding)
	var views []*CompiledView
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		if path.Ext(file.Name()) != ".svelte" {
			continue
		}
		view := strings.TrimSuffix(file.Name(), ".svelte")
		dom, err := c.compileDOM(view)
		if err != nil {
			return nil, err
		}
		h := sha1.New()
		h.Write(dom)
		domHash := encoding.EncodeToString(h.Sum(nil))
		ssr, err := c.compileSSR(view)
		if err != nil {
			return nil, err
		}
		h = sha1.New()
		h.Write(ssr)
		ssrHash := encoding.EncodeToString(h.Sum(nil))
		css, err := c.extractCSS(ssr)
		if err != nil {
			return nil, err
		}
		h = sha1.New()
		h.Write(css)
		cssHash := encoding.EncodeToString(h.Sum(nil))
		views = append(views, &CompiledView{
			Name:    view,
			CSS:     css,
			CSSHash: strings.ToLower(cssHash),
			DOM:     dom,
			DOMHash: strings.ToLower(domHash),
			SSR:     ssr,
			SSRHash: strings.ToLower(ssrHash),
		})

	}
	return views, nil
}

type SSRResult struct {
	HTML string `json:"html"`
	CSS  SSRCSS `json:"css"`
	HEAD string `json:"head"`
}

type SSRCSS struct {
	Code string `json:"code"`
}

func (c *Compiler) extractCSS(ssrScript []byte) ([]byte, error) {
	ssr, err := ssr.Load(c.vm, string(ssrScript))
	if err != nil {
		return nil, err
	}
	var dst SSRResult
	err = ssr.Execute(&dst, map[string]interface{}{})
	if err != nil {
		return nil, err
	}
	return []byte(dst.CSS.Code), nil
}

func (c *Compiler) compileDOM(view string) ([]byte, error) {
	res := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{path.Join(c.dir, "views", fmt.Sprintf("%s.svelte", view))},
		Bundle:            true,
		Outfile:           path.Join("dom.js"),
		Format:            esbuild.FormatIIFE,
		GlobalName:        "__dom__",
		MinifyWhitespace:  true,
		IgnoreAnnotations: true,
		Sourcemap:         esbuild.SourceMapNone,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Platform:          esbuild.PlatformBrowser,
		Plugins:           c.tm.DOM.Plugins(),
		AbsWorkingDir:     c.wd,
	})
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("wx: failed to compile DOM %s: %s", view, res.Errors[0].Text)
	}
	return res.OutputFiles[0].Contents, nil
}

func (c *Compiler) compileSSR(view string) ([]byte, error) {
	res := esbuild.Build(esbuild.BuildOptions{
		EntryPoints:       []string{path.Join(c.dir, "views", fmt.Sprintf("%s.svelte", view))},
		Bundle:            true,
		Outfile:           path.Join("ssr.js"),
		Format:            esbuild.FormatIIFE,
		GlobalName:        "__ssr__",
		MinifyWhitespace:  true,
		IgnoreAnnotations: true,
		Sourcemap:         esbuild.SourceMapNone,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Platform:          esbuild.PlatformNode,
		Plugins:           c.tm.SSR.Plugins(),
		AbsWorkingDir:     c.wd,
	})
	if len(res.Errors) > 0 {
		return nil, fmt.Errorf("wx: failed to compile SSR %s: %s", view, res.Errors[0].Text)
	}
	return res.OutputFiles[0].Contents, nil
}
