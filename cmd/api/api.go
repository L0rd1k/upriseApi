package api

import "net/http"

type Api struct {
	stack []Middleware
	app   App
}

func NewApi() *Api {
	return &Api{
		stack: []Middleware{},
		app:   nil,
	}
}

func (api *Api) Use(middlewares ...Middleware) {
	api.stack = append(api.stack, middlewares...)
}

func (api *Api) SetApp(app App) {
	api.app = app
}

func (api *Api) MakeHandler() http.Handler {
	var appFunc HandlerFunc
	if api.app != nil {
		appFunc = api.app.AppFunc()
	} else {
		appFunc = func(w ResponseWriter, r *Request) {}
	}
	return http.HandlerFunc(
		adapterFunc(
			WrapMiddlewares(api.stack, appFunc),
		),
	)
}
