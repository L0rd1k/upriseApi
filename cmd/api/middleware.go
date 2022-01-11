package api

import "net/http"

type HandlerFunc func(ResponseWriter, *Request)
type MiddlewareSimple func(handler HandlerFunc) HandlerFunc
type AppSimple HandlerFunc

type Middleware interface {
	MiddlewareFunc(handler HandlerFunc) HandlerFunc
}

type App interface {
	AppFunc() HandlerFunc
}

func (m_ware MiddlewareSimple) MiddlewareFunc(handler HandlerFunc) HandlerFunc {
	return m_ware(handler)
}

func (a_smpl AppSimple) AppFunc() HandlerFunc {
	return HandlerFunc(a_smpl)
}

func WrapMiddlewares(middlewares []Middleware, handler HandlerFunc) HandlerFunc {
	wrapped := handler
	for i := len(middlewares) - 1; i >= 0; i-- {
		wrapped = middlewares[i].MiddlewareFunc(wrapped)
	}
	return wrapped
}

func adapterFunc(handler HandlerFunc) http.HandlerFunc {
	return func(origWriter http.ResponseWriter, origRequest *http.Request) {
		request := &Request{
			origRequest,
			nil,
			map[string]interface{}{},
		}

		writer := &responseWriter{
			origWriter,
			false,
		}
		handler(writer, request)
	}
}
