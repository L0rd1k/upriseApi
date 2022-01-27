package router

import (
	"strings"

	"github.com/L0rd1k/uprise-api/cmd/api"
)

type Route struct {
	HttpMethod     string
	PathExpression string
	Function       api.HandlerFunc
}

func (route *Route) MakePath(pathParams map[string]string) string {
	path := route.PathExpression
	for paramName, paramValue := range pathParams {
		paramPlaceholder := ":" + paramName
		relaxPlaceholder := "#" + paramName
		splatPlaceholder := "*" + paramName
		rqst := strings.NewReplacer(
			paramPlaceholder, paramValue,
			splatPlaceholder, paramValue,
			relaxPlaceholder, paramValue)
		path = rqst.Replace(path)
	}
	return path
}

func Head(pathExp string, HandlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "HEAD",
		PathExpression: pathExp,
		Function:       HandlerFunc,
	}
}

func Get(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "GET",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}

func Post(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "POST",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}

func Put(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "PUT",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}

func Patch(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "PATCH",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}

func Delete(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "DELETE",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}

func Options(pathExp string, handlerFunc api.HandlerFunc) *Route {
	return &Route{
		HttpMethod:     "OPTIONS",
		PathExpression: pathExp,
		Function:       handlerFunc,
	}
}
