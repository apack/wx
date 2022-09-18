package wx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
)

type Context struct {
	w   http.ResponseWriter
	r   *http.Request
	eng *RenderEngine
	context.Context
}

func (ctx *Context) Method() string {
	return ctx.r.Method
}

func (ctx *Context) Path() string {
	return ctx.r.URL.Path
}

func (ctx *Context) Query() url.Values {
	return ctx.r.URL.Query()
}

func (ctx *Context) ParseForm() (url.Values, error) {
	err := ctx.r.ParseForm()
	if err != nil {
		return nil, err
	}
	return ctx.r.Form, nil
}

func (ctx *Context) JSON(v interface{}) error {
	ctx.w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(ctx.w)
	return enc.Encode(v)
}

func (ctx *Context) View(view View, props Props) error {
	html, err := ctx.eng.Render(view, props)
	if err != nil {
		return err
	}
	ctx.w.Header().Set("Content-Type", "text/html")
	_, err = ctx.w.Write(html)
	return err
}

func (ctx *Context) Redirect(url string, code int) error {
	http.Redirect(ctx.w, ctx.r, url, code)
	return nil
}

func (ctx *Context) Status(code int) *Context {
	ctx.w.WriteHeader(code)
	return ctx
}
