package main

import (
	"net/http"
	"strings"
)

type HandlerBaseOnMap struct {
	handlers map[string]func(ctx *Context)
}

func (h *HandlerBaseOnMap) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	key := h.key(r.Method, r.URL.Path)
	if handler, ok := h.handlers[key]; ok {
		handler(NewContext(w, r))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}
}

func (h *HandlerBaseOnMap) key(method, pattern string) string {
	return strings.ToUpper(method) + "#" + pattern
}
