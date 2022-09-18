package main

import (
	"net/http"

	"github.com/apack/wx"
	"github.com/apack/wx/example/app"
)

func main() {
	web := wx.NewRouter()
	err := app.Load(web)
	if err != nil {
		panic(err)
	}
	web.HandleMethodFunc("GET", "/", handleWelcomeView())
	http.ListenAndServe(":8080", web)
}

func handleWelcomeView() wx.HandlerFunc {
	return func(ctx *wx.Context) error {
		return ctx.View(app.WelcomeView, wx.Props{
			"count": 10,
		})
	}
}
