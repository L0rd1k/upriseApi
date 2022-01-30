package router

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/L0rd1k/uprise-api/cmd/api"
	"github.com/L0rd1k/uprise-api/cmd/api/trie"
)

var preEscape = strings.NewReplacer("*", "__SPLAT_PLACEHOLDER__", "#", "__RELAXED_PLACEHOLDER__")
var postEscape = strings.NewReplacer("__SPLAT_PLACEHOLDER__", "*", "__RELAXED_PLACEHOLDER__", "#")

type Router struct {
	Routes         []*Route
	index          map[*Route]int
	compressionOff bool
	trie           *trie.Trie
}

func (router *Router) AppFunc() api.HandlerFunc {
	return func(writer api.ResponseWriter, request *api.Request) {
		route, params, pathMatched := router.findRouteFromURL(request.Method, request.URL)
		if route == nil {
			if pathMatched {
				Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
				return
			}
			NotFound(writer, request)
			return
		}
		request.PathParams = params
		handler := route.Function
		handler(writer, request)
	}
}

func MakeRouter(routes ...*Route) (api.App, error) {
	r := &Router{
		Routes: routes,
	}
	error := r.start()
	if error != nil {
		return nil, error
	}
	return r, nil
}

func (router *Router) findRouteFromURL(httpMethod string, urlObj *url.URL) (*Route, map[string]string, bool) {
	matches, pathMatched := router.trie.FindRoutesAndPathMatched(
		strings.ToUpper(httpMethod),
		escapedPath(urlObj),
	)
	if len(matches) == 0 {
		return nil, nil, pathMatched
	}

	if len(matches) == 1 {
		return matches[0].Route.(*Route), matches[0].Params, pathMatched
	}
	result := router.ofFirstDefinedRoute(matches)
	return result.Route.(*Route), result.Params, pathMatched
}

func (router *Router) start() error {
	router.trie = trie.New()
	router.index = map[*Route]int{}
	for i, route := range router.Routes {
		pathExp, error := escapedPathExp(route.PathExpression)
		if error != nil {
			return error
		}
		error = router.trie.AddRoute(
			strings.ToUpper(route.HttpMethod),
			pathExp,
			route,
		)
		if error != nil {
			return error
		}
		router.index[route] = i
	}

	if router.compressionOff == false {
		router.trie.Compress()
	}
	return nil
}

func (router *Router) ofFirstDefinedRoute(matches []*trie.Match) *trie.Match {
	minIndex := -1
	var bestMatch *trie.Match
	for _, result := range matches {
		route := result.Route.(*Route)
		routeIndex := router.index[route]
		if minIndex == -1 || routeIndex < minIndex {
			minIndex = routeIndex
			bestMatch = result
		}
	}
	return bestMatch
}

func escapedPathExp(pathExp string) (string, error) {
	if pathExp == "" {
		return "", errors.New("empty path expression")
	} else if pathExp[0] != '/' {
		return "", errors.New("path expression must start with /")
	} else if strings.Contains(pathExp, "?") {
		return "", errors.New("path expression must not contain the query string")
	}
	pathExp = preEscape.Replace(pathExp)
	urlObj, error := url.Parse(pathExp)
	if error != nil {
		return "", error
	}
	pathExp = urlObj.RequestURI()
	pathExp = postEscape.Replace(pathExp)
	return pathExp, nil
}

func escapedPath(urlObj *url.URL) string {
	parts := strings.SplitN(urlObj.RequestURI(), "?", 2)
	return parts[0]
}

func Error(w api.ResponseWriter, error string, code int) {
	w.WriteHeader(code)
	err := w.WriteJson(map[string]string{api.ErrorFieldName: error})
	if err != nil {
		panic(err)
	}
}

func NotFound(w api.ResponseWriter, r *api.Request) {
	Error(w, "Resource not found", http.StatusNotFound)
}
