package main

import (
	"net/http"
	"sync"

	"github.com/apack/wx"
	"github.com/apack/wx/example/app"
)

func main() {
	web := wx.NewRouter()
	err := app.Load(web)
	if err != nil {
		panic(err)
	}
	store := new(Store)
	web.HandleMethodFunc("GET", "/", handleWelcomeView(store))
	web.HandleMethodFunc("POST", "/api/increase", handleIncrease(store))
	web.HandleMethodFunc("POST", "/api/decrease", handleDecrease(store))
	web.HandleMethodFunc("GET", "/api/count", handleCount(store))
	http.ListenAndServe(":8080", web)
}

func handleWelcomeView(store *Store) wx.HandlerFunc {
	return func(ctx *wx.Context) error {
		return ctx.View(app.WelcomeView, wx.Props{
			"count": store.GetCount(),
		})
	}
}

type Store struct {
	Count int
	mtx   sync.RWMutex
}

func (s *Store) Increase() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Count++
	return s.Count
}

func (s *Store) GetCount() int {
	s.mtx.RLock()
	defer s.mtx.RUnlock()
	return s.Count
}

func (s *Store) Decrease() int {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.Count--
	return s.Count
}

func handleIncrease(store *Store) wx.HandlerFunc {
	return func(ctx *wx.Context) error {
		return ctx.JSON(wx.Props{
			"count": store.Increase(),
		})
	}
}

func handleDecrease(store *Store) wx.HandlerFunc {
	return func(ctx *wx.Context) error {
		return ctx.JSON(wx.Props{
			"count": store.Decrease(),
		})
	}
}

func handleCount(store *Store) wx.HandlerFunc {
	return func(ctx *wx.Context) error {
		return ctx.JSON(wx.Props{
			"count": store.GetCount(),
		})
	}
}
