package api

import (
	"fmt"
	"net/http"
)

type Api struct {
	stack []Middleware //> list of middleware
	app   App          //> App object
}

func NewApi() *Api {
	return &Api{
		stack: []Middleware{},
		app:   nil,
	}
}

func (api *Api) Use(middlewares ...Middleware) {
	fmt.Println("<Api:Use> api.stack:", api.stack)
	api.stack = append(api.stack, middlewares...)
	fmt.Println("<Api:Use> api.stack:", api.stack)
}

func (api *Api) SetApp(app App) {
	fmt.Println("<Api:SetApp>")
	api.app = app
}

func (api *Api) MakeHandler() http.Handler {
	var appFunc HandlerFunc
	if api.app != nil {
		fmt.Println("<Api:MakeHandler>App not null")
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

var DefaultCommonStack = []Middleware{
	&TimerMiddleware{},
	&RecorderMiddleware{},
}
