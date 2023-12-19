package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"time"
)

// 定义中间
func MiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		fmt.Println("中间件开始执行了")
		// 设置变量到Context的key中，可以通过Get()取
		c.Set("request", "中间件")
		c.Next()
		status := c.Writer.Status()
		fmt.Println("中间件执行完毕", status)
		t2 := time.Since(t)
		fmt.Println("time:", t2)
	}
}

func MiddleWare2() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("中间件2开始执行了")
		c.Set("request2", "中间件2")
		c.Next()
		fmt.Println("中间件2执行完毕")
	}
}

func main() {
	// 1.创建路由
	// 默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()

	// 注册中间件
	r.Use(MiddleWare())
	r.Use(MiddleWare2())

	r.GET("/ce", func(c *gin.Context) {
		// 取值
		req, _ := c.Get("request")
		req2, _ := c.Get("request2")
		fmt.Println("request:", req)
		// 页面接收
		c.JSON(200, gin.H{"request": req, "request2:": req2})
	})

	r.Run(":8000")
}
