package main

import (
	"net/http"
)

type Routable interface {
	Route(method, pattern string, handlerFunc func(ctx *Context))
}

type Handler interface {
	ServeHTTP(ctx *Context)
	Routable
}

type HandlerBaseOnMap struct {
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBaseOnMap) ServeHTTP(ctx *Context) {
	key := h.key(ctx.R.Method, ctx.R.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(ctx)
	} else {
		ctx.W.WriteHeader(http.StatusNotFound)
		ctx.W.Write([]byte("Not Found"))
		return
	}
}

// Route 路由注册
func (h *HandlerBaseOnMap) Route(method, pattern string, handlerFunc func(ctx *Context)) {
	key := h.key(method, pattern)
	h.handlers[key] = handlerFunc
}

func (h *HandlerBaseOnMap) key(method, pattern string) string {
	return method + "#" + pattern
}

// 用于确保 HandlerBaseOnMap 实现了 Handler 接口
var _ Handler = &HandlerBaseOnMap{}

func NewHandlerBaseOnMap() Handler {
	return &HandlerBaseOnMap{
		make(map[string]func(ctx *Context)),
	}
}
