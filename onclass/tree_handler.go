package main

import (
	"net/http"
	"strings"
)

type HandlerBaseOnTree struct {
	root *Node
}

type Node struct {
	path     string
	children []*Node
	handler  handleFunc
}

type handleFunc func(ctx *Context)

func (h *HandlerBaseOnTree) Route(method, pattern string, handlerFunc handleFunc) {
	// 去掉前后的 '/'
	pattern = strings.Trim(pattern, "/")
	paths := strings.Split(pattern, "/")

	// 指向当前根节点
	cur := h.root
	for index, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if ok {
			cur = matchChild
		} else {
			cur.createChild(paths[index:], handlerFunc)
			break
		}
	}
	cur.handler = handlerFunc
}

func (h *HandlerBaseOnTree) Start(address string) error {
	return nil
}

func (h *HandlerBaseOnTree) ServeHTTP(ctx *Context) {
	handler, ok := h.findRouter(ctx.R.URL.Path)
	if !ok {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("Not Found"))
		return
	}
	handler(ctx)
}

func (h *HandlerBaseOnTree) findRouter(path string) (handleFunc, bool) {
	paths := strings.Split(strings.Trim(path, "/"), "/")
	cur := h.root
	for _, path := range paths {
		matchChild, ok := cur.findMatchChild(path)
		if !ok {
			return nil, false
		}
		cur = matchChild
	}

	if cur.handler == nil {
		return nil, false
	}
	return cur.handler, true
}

func (n *Node) findMatchChild(path string) (*Node, bool) {
	var wildcardNode *Node
	for _, child := range n.children {
		// 严格匹配
		if child.path == path && child.path != "*" {
			return child, true
		}
		// 通配符匹配
		if child.path == "*" {
			wildcardNode = child
		}
	}
	return wildcardNode, wildcardNode != nil
}

func (n *Node) createChild(paths []string, handlerFunc handleFunc) {
	for _, path := range paths {
		nn := newNode(path)
		n.children = append(n.children, nn)
		n = nn
	}
}

func newNode(path string) *Node {
	return &Node{
		path:     path,
		children: make([]*Node, 0, 2),
	}
}
