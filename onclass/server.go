package main

import (
	"net/http"
)

// Server http server 的顶级抽象
type Server interface {
	Routable
	Start(address string) error
}

// sdkHttpServer 基于 http 库实现
type sdkHttpServer struct {
	Name    string
	handler Handler
	root    Filter
}

// Route 路由注册
func (s sdkHttpServer) Route(method, pattern string, handlerFunc func(ctx *Context)) {
	s.handler.Route(method, pattern, handlerFunc)
}

// Start 启动服务
func (s sdkHttpServer) Start(address string) error {
	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		ctx := NewContext(writer, request)
		s.root(ctx)
	})
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 创建服务
func NewHttpServer(name string, builders ...FilterBuilder) Server {
	handler := NewHandlerBaseOnMap()
	var root Filter = func(ctx *Context) {
		handler.ServeHTTP(ctx)
	}

	for i := len(builders) - 1; i >= 0; i-- {
		b := builders[i]
		root = b(root)
	}

	return &sdkHttpServer{
		Name:    name,
		handler: handler,
		root:    root,
	}
}

func SignUp(ctx *Context) {

	req := &SignUpReq{}

	err := ctx.ReadJson(req)
	if err != nil {
		ctx.BadRequestJson(err)
		return
	}

	data := &CommonResp{
		Data: "123",
	}
	err = ctx.OkJson(data)
	if err != nil {
		ctx.SystemErrorJson(err)
		return
	}

}

type SignUpReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type CommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}
