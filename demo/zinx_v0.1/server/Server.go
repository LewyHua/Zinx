package main

import (
	"zinx/znet"
)

// 基于zinx框架开发的 服务器端应用程序

func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer("[zinx v0.1]")
	s.Serve()
}
