package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("welcome home"))
}

func wholeUrl(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(r.URL)
	fmt.Fprintf(w, string(data))
}

func user(w http.ResponseWriter, r *http.Request) {

}

func main() {
	server := NewHttpServer("my-server")

	//server.Route("/", home)
	//server.Route("/wholeUrl", wholeUrl)
	// 注册路由
	server.Route("post", "/signUp", SignUp)

	// 启动服务
	err := server.Start("localhost:9090")
	if err != nil {
		fmt.Printf("start server failed, err: %v\n", err)
	}

}
