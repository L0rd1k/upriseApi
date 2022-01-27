package router

import (
	"errors"
	"net/url"
	"strings"

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

// func MakeRouter(routes ...*Route) (App, error) {
// 	r := &Router{
// 		Routes: routes,
// 	}
// 	r.start()
// }

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

func escapedPathExp(pathExp string) (string, error) {
	if pathExp == "" {
		return "", errors.New("Empty path expression")
	} else if pathExp[0] != '/' {
		return "", errors.New("Path expression must start with /")
	} else if strings.Contains(pathExp, "?") {
		return "", errors.New("Path expression must not contain the query string")
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
