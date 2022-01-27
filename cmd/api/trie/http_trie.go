package trie

import (
	"errors"
	"fmt"
)

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

type Trie struct {
	root *Node
}

func New() *Trie {
	return &Trie{
		root: &Node{},
	}
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

//--------------------------------------------------------------------------------------------

func (t *Trie) AddRoute(httpMethod, pathExp string, route interface{}) error {
	return t.root.addRoute(httpMethod, pathExp, route, []string{})
}

func (t *Trie) Compress() {
	t.root.compress()
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
