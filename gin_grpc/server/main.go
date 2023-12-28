package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	pb "go-web-study/gin_grpc/server/proto"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type server struct {
	pb.UnimplementedSayHelloServer
}

func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	fmt.Printf("进入到 grpc 实现\n")
	return &pb.HelloResponse{ResponseMsg: "hello" + req.RequestName}, nil
}

// 入口
//func main() {
//	// 初始化grpc服务
//	grpcServer := grpc.NewServer()
//
//	/***** 注册你的grpc服务 *****/
//	pb.RegisterSayHelloServer(grpcServer, &server{})
//
//	// 初始化一个空Gin路由
//	router := gin.New()
//
//	/***** 添加你的api路由吧 *****/
//	router.GET("/test", func(ctx *gin.Context) {
//		fmt.Printf("已进入到HTTP实现\n")
//		ctx.JSON(200, map[string]interface{}{
//			"code": 200,
//			"msg:": "test http success",
//		})
//	})
//
//	// 监听端口并处理服务分流
//	h2Handler := h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
//		// 判断协议是否为http/2 && 是grpc
//		if r.ProtoMajor == 2 &&
//			strings.HasPrefix(r.Header.Get("Content-Type"), "application/grpc") {
//			// 按grpc方式来请求
//			fmt.Printf("按grpc方式来请求\n")
//			grpcServer.ServeHTTP(w, r)
//		} else {
//			// 当作普通api
//			fmt.Printf("当作普通api\n")
//			router.ServeHTTP(w, r)
//		}
//	}), &http2.Server{})
//
//	// 监听HTTP服务
//	if err := http.ListenAndServe(":8080", h2Handler); err != nil {
//		log.Println("http server done:", err.Error())
//	}
//}

// 自定义的 ResponseWriter 类型
type responseWriter struct {
	http.ResponseWriter
}

// 实现 http.ResponseWriter 接口的 Write 方法
func (rw *responseWriter) Write(data []byte) (int, error) {
	// 在写入响应之前可以执行一些额外的操作
	fmt.Println("Writing response:", string(data))

	// 调用原始的 Write 方法
	return rw.ResponseWriter.Write(data)
}

// 入口
func main() {
	// 初始化grpc服务
	grpcServer := grpc.NewServer()

	/***** 注册你的grpc服务 *****/
	pb.RegisterSayHelloServer(grpcServer, &server{})

	// 初始化一个空Gin路由
	router := gin.New()

	// 全局拦截
	router.Use(func(ctx *gin.Context) {
		// 判断协议是否为http/2
		// 判断是否是grpc
		if ctx.Request.ProtoMajor == 2 &&
			strings.HasPrefix(ctx.Request.Header.Get("Content-Type"), "application/grpc") {
			// 按grpc方式来请求
			fmt.Printf("按grpc方式来请求\n")
			rw := &responseWriter{ResponseWriter: ctx.Writer}
			grpcServer.ServeHTTP(rw, ctx.Request)
			// 不要再往下请求了,防止继续链式调用拦截器
			ctx.Abort()
			return
		}
		// 当作普通api
		fmt.Printf("当作普通api\n")
		ctx.Next()
	})

	/***** 添加你的api路由吧 *****/
	router.GET("/test", func(c *gin.Context) {
		fmt.Printf("进入到 http 实现\n")
		c.JSON(200, map[string]interface{}{
			"code": 200,
			"msg:": "test http success",
		})
	})

	// 为http/2配置参数
	h2Handle := h2c.NewHandler(router, &http2.Server{}) // 禁用TLS加密协议
	// 配置http服务
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2Handle,
	}
	// 启动http服务
	server.ListenAndServe()
}

//func main() {
//	// 开启服务，监听端口
//	listen, _ := net.Listen("tcp", ":8080")
//
//	// 创建 grpc 服务
//	grpcServer := grpc.NewServer()
//
//	// 注册服务
//	pb.RegisterSayHelloServer(grpcServer, &server{})
//
//	// 启动服务
//	err := grpcServer.Serve(listen)
//	if err != nil {
//		fmt.Printf("failed to server, err: %v\n", err)
//		return
//	}
//	fmt.Printf("server start success...\n")
//}
