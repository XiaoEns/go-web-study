package main

import (
	"context"
	"fmt"
	pb "go-web-study/gin_grpc/server/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net/http"
)

//func main() {
//	var Address = os.Args[1]
//	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stdout, os.Stdout))
//
//	// 连接
//	conn, err := grpc.Dial(Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	if err != nil {
//		grpclog.Fatalln(err)
//	}
//	defer conn.Close()
//
//	// 初始化客户端
//	c := pb.NewSayHelloClient(conn)
//
//	// 调用方法
//	resp, err := c.SayHello(context.Background(), &pb.HelloRequest{RequestName: " 肖恩"})
//	if err != nil {
//		grpclog.Fatalln(err)
//	}
//
//	fmt.Printf("resp = %v\n", resp)
//}

func main() {
	conn, err := grpc.Dial("127.0.0.1:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("connected failed, err:%v\n", err)
	}
	defer conn.Close()

	// 建立连接
	client := pb.NewSayHelloClient(conn)

	// 执行 rpc 调用
	resp, err := client.SayHello(context.Background(), &pb.HelloRequest{RequestName: " 肖恩"})
	if err != nil {
		fmt.Printf("grpc failed, err: %v\n", err)
	} else {
		fmt.Printf("grpc resp = %v\n", resp)
	}

	// 执行 http 调用
	resp2, err := http.NewRequest("GET", "http://localhost:8080/test", nil)
	if err != nil {
		fmt.Printf("http failed, err: %v\n", err)
	} else {
		fmt.Printf("http resp = %v\n", resp2)
	}

}
