![wx](/assets/logo.png)

# wx


wx is an open-source framework that allows you to quickly build Go apps that use Svelte for rendering web pages.

This project was heavily inspired by [bud](https://github.com/livebud/bud). wx tries to be less opinionated about your Go app build process.


## IMPORTANT!

This project is at the very early stage. Do not use it for production! Wait until **`v1`** is released. 

## Usage

### Installation

```sh
go get github.com/apack/wx
go install github.com/apack/wx/cmd/...
```

### Initialize Views

Open your Go app folder and initialize views.
```sh
wx init
```
This will create files:
```txt
app
├── app.wx.go
├── components
│   ├── Button.svelte
│   └── Counter.svelte
├── gen.go
├── layouts
│   └── DefaultLayout.svelte
├── static
│   └── logo.svg
└── views
    └── WelcomeView.svelte
```

**NOTE:** Every time you change these files run `go generate ./app`. 
### Create Server

Load your app and add handler.

```go
package main

import (
	"net/http"
	"your/path/to/app"

	"github.com/apack/wx"
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
			"count": 111,
		})
	}
}

```

## Roadmap

These are the main accomplishments that I want to release with future versions.

### `v0.1` - Done
* ~~Create concept and release it to the public~~

### `v0.2` - WIP
* Improve view renderer.

### `v0.3`
* Improve router.
* Add middlewares for CORS, Caching and Compression.

### `v0.4`
* Make project ready for community contributions.

### `v0.5`
* Add TypeScript support for Svelte components.
* Add SASS support for Svelte components.
* Add Markdown support for static content.

### `v0.6`
* Fix memory issue with long-lived v8 isolates. [v8go#105](https://github.com/rogchap/v8go/issues/105)


### `v0.7`
* Create documentation and landing pages using `wx`.


## Contributing

Project is not ready to accept contributors yet. Wait until `v0.4` when all of the development documentation and CI is finished.
## Credits

This work is based off the existing frameworks:
* [bud](https://github.com/livebud/bud)

## License

Copyright (c) 2022-present APack and Contributors. wx is free and open-source software licensed under the MIT License. 

Third-party library licenses:

* [v8go](https://github.com/rogchap/v8go/blob/master/LICENSE)
* [v8go-polyfills](https://github.com/kuoruan/v8go-polyfills/blob/master/LICENSE)
* [Svelte](https://github.com/sveltejs/svelte/blob/master/LICENSE.md)
* [esbuild](https://github.com/evanw/esbuild/blob/master/LICENSE.md)
* [cobra](https://github.com/spf13/cobra/blob/main/LICENSE.txt)