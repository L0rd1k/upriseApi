package trie

import (
	"errors"
	"fmt"
)

type Match struct {
	Route  interface{}
	Params map[string]string
}

type paramMatch struct {
	name  string
	value string
}

func New() *Trie {
	return &Trie{
		root: &Node{},
	}
}

func newFindContext() *findContext {
	return &findContext{
		paramStack: []paramMatch{},
	}
}

//--------------------------------------------------------------------------------------------

type findContext struct {
	paramStack []paramMatch
	matchFunc  func(httpMethod, path string, node *Node)
}

func (fc *findContext) paramsAsMap() map[string]string {
	r := map[string]string{}
	for _, param := range fc.paramStack {
		if r[param.name] != "" {
			panic(fmt.Sprintf(
				"Placeholder %s already found, placeholder names should be unique per route", param.name,
			))
		}
		r[param.name] = param.value
	}
	return r
}

func (fc *findContext) pushParams(name, value string) {
	fc.paramStack = append(fc.paramStack, paramMatch{name, value})
}

func (fc *findContext) popParams() {
	fc.paramStack = fc.paramStack[:len(fc.paramStack)-1]
}

//--------------------------------------------------------------------------------------------

type Node struct {
	HttpMethodToRoute map[string]interface{}
	Children          map[string]*Node
	ParamChild        *Node
	RelaxedChild      *Node
	SplatChild        *Node

	ChildrenKeyLen int
	ParamName      string
	RelaxedName    string
	SplatName      string
}

func (nd *Node) addRoute(httpMethod, pathExp string, route interface{}, usedParams []string) error {
	if len(pathExp) == 0 {
		if nd.HttpMethodToRoute == nil {
			nd.HttpMethodToRoute = map[string]interface{}{
				httpMethod: route,
			}
			return nil
		} else {
			if nd.HttpMethodToRoute[httpMethod] != nil {
				return errors.New("Node.route already set, duplicated path and method")
			}
			nd.HttpMethodToRoute[httpMethod] = route
			return nil
		}
	}
	token := pathExp[0:1]
	remaining := pathExp[1:]
	var nextNode *Node
	switch token[0] {
	case ':':
		var name string
		name, remaining = splitParam(remaining)
		for _, elem := range usedParams {
			if elem == name {
				return errors.New(fmt.Sprintf("A route can't have two placeholders with the same name: %s", name))
			}
		}
		usedParams = append(usedParams, name)
		if nd.ParamChild == nil {
			nd.ParamChild = &Node{}
			nd.ParamName = name
		} else {
			if nd.ParamName != name {
				return errors.New(
					fmt.Sprintf("Routes sharing a common placeholder MUST name it consistently: %s != %s",
						nd.ParamName, name),
				)
			}
		}
		nextNode = nd.ParamChild
	case '#':
		var name string
		name, remaining = splitRelaxed(remaining)
		for _, elem := range usedParams {
			if elem == name {
				return errors.New(fmt.Sprintf("A route can't have two placeholders with the same name: %s", name))
			}
		}
		usedParams = append(usedParams, name)
		if nd.RelaxedChild == nil {
			nd.RelaxedChild = &Node{}
			nd.RelaxedName = name
		} else {
			if nd.RelaxedName != name {
				return errors.New(
					fmt.Sprintf("Routes sharing a common placeholder MUST name it consistently: %s != %s",
						nd.RelaxedName, name),
				)
			}
		}
		nextNode = nd.RelaxedChild
	case '*':
		name := remaining
		remaining = ""
		for _, elem := range usedParams {
			if elem == name {
				return errors.New(fmt.Sprintf("A route can't have two placeholders with the same name: %s", name))
			}
		}
		if nd.SplatChild == nil {
			nd.SplatChild = &Node{}
			nd.SplatName = name
		}
		nextNode = nd.SplatChild
	default:
		if nd.Children == nil {
			nd.Children = map[string]*Node{}
			nd.ChildrenKeyLen = 1
		}
		if nd.Children[token] == nil {
			nd.Children[token] = &Node{}
		}
		nextNode = nd.Children[token]
	}
	return nextNode.addRoute(httpMethod, remaining, route, usedParams)
}

func (nd *Node) compress() {
	if nd.SplatChild != nil {
		nd.SplatChild.compress()
	}
	if nd.ParamChild != nil {
		nd.ParamChild.compress()
	}
	if nd.RelaxedChild != nil {
		nd.RelaxedChild.compress()
	}
	if len(nd.Children) == 0 {
		return
	}
	canCompress := true
	for _, node := range nd.Children {
		if node.HttpMethodToRoute != nil || node.SplatChild != nil ||
			node.ParamChild != nil || node.RelaxedChild != nil {
			canCompress = false
		}
	}

	if canCompress {
		merged := map[string]*Node{}
		for key, node := range nd.Children {
			for gdKey, gdNode := range node.Children {
				mergedKey := key + gdKey
				merged[mergedKey] = gdNode
			}
		}
		nd.Children = merged
		nd.ChildrenKeyLen++
		nd.compress()
	} else {
		for _, node := range nd.Children {
			node.compress()
		}
	}
}

func (nd *Node) find(httpMethod, path string, context *findContext) {
	if nd.HttpMethodToRoute != nil && path == "" {
		context.matchFunc(httpMethod, path, nd)
	}
	if len(path) == 0 {
		return
	}

	if nd.SplatChild != nil {
		context.pushParams(nd.SplatName, path)
		nd.SplatChild.find(httpMethod, "", context)
		context.popParams()
	}

	if nd.ParamChild != nil {
		value, remaining := splitParam(path)
		context.pushParams(nd.ParamName, value)
		nd.ParamChild.find(httpMethod, remaining, context)
		context.popParams()
	}

	if nd.RelaxedChild != nil {
		value, remaining := splitRelaxed(path)
		context.pushParams(nd.RelaxedName, value)
		nd.RelaxedChild.find(httpMethod, remaining, context)
		context.popParams()
	}
	length := nd.ChildrenKeyLen
	if len(path) < length {
		return
	}
	token := path[0:length]
	remaining := path[length:]
	if nd.Children[token] != nil {
		nd.Children[token].find(httpMethod, remaining, context)
	}
}

//--------------------------------------------------------------------------------------------

type Trie struct {
	root *Node
}

func (t *Trie) AddRoute(httpMethod, pathExp string, route interface{}) error {
	return t.root.addRoute(httpMethod, pathExp, route, []string{})
}

func (t *Trie) Compress() {
	t.root.compress()
}

func (t *Trie) FindRoutesAndPathMatched(httpMethod, path string) ([]*Match, bool) {
	context := newFindContext()
	pathMatched := false
	matches := []*Match{}
	context.matchFunc = func(httpMethod, path string, node *Node) {
		pathMatched = true
		if node.HttpMethodToRoute[httpMethod] != nil {
			matches = append(matches, &Match{
				Route:  node.HttpMethodToRoute[httpMethod],
				Params: context.paramsAsMap(),
			})
		}
	}
	t.root.find(httpMethod, path, context)
	return matches, pathMatched
}

//--------------------------------------------------------------------------------------------

func splitParam(remaining string) (string, string) {
	i := 0
	for len(remaining) > i && remaining[i] != '/' && remaining[i] != '.' {
		i++
	}
	return remaining[:i], remaining[i:]
}

func splitRelaxed(remaining string) (string, string) {
	i := 0
	for len(remaining) > i && remaining[i] != '/' {
		i++
	}
	return remaining[:i], remaining[i:]
}
