package main

import (
	"net/http"
)

// Server http server 的顶级抽象
type Server interface {
	Route(method, pattern string, handlerFunc func(ctx *Context))
	Start(address string) error
}

// sdkHttpServer 基于 http 库实现
type sdkHttpServer struct {
	Name    string
	handler *HandlerBaseOnMap
}

// Route 路由注册
func (s sdkHttpServer) Route(method, pattern string, handlerFunc func(ctx *Context)) {
	key := s.handler.key(method, pattern)
	s.handler.handlers[key] = handlerFunc
}

// Start 启动服务
func (s sdkHttpServer) Start(address string) error {
	http.Handle("/", s.handler)
	return http.ListenAndServe(address, nil)
}

// NewHttpServer 创建服务
func NewHttpServer(name string) Server {
	return &sdkHttpServer{
		Name:    name,
		handler: &HandlerBaseOnMap{make(map[string]func(ctx *Context))},
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
