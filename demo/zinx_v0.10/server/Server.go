package main

import (
	"fmt"
	"sync/atomic"
	"zinx/ziface"
	"zinx/znet"
)

// 基于zinx框架开发的 服务器端应用程序
func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer(znet.WithName("HAHAHA"))

	// 2 注册开启关闭conn时的hook函数
	s.RegisterOnConnStart(StartHookFunc)
	s.RegisterOnConnStop(StopHookFunc)

	// 3 给当前Zinx框架增加自定义Router
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	s.Serve()
}

// PingRouter test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
	var atomicInt atomic.Uint32
	fmt.Println("Calling PingRouter...")
	fmt.Printf("MsgID: %d, Data: %v\n", request.GetMsgID(), string(request.GetMsgData()))
	err := request.GetConn().SendMsg(atomicInt.Add(200), []byte("Pong...\n"))
	if err != nil {
		fmt.Println("Handle err:", err)
		return
	}
}

// HelloRouter test自定义路由
type HelloRouter struct {
	znet.BaseRouter
}

func (p *HelloRouter) Handle(request ziface.IRequest) {
	var atomicInt atomic.Uint32
	//fmt.Println("Calling HelloHandler!!!")
	fmt.Printf("MsgID: %d, Data: %v\n", request.GetMsgID(), string(request.GetMsgData()))
	err := request.GetConn().SendMsg(atomicInt.Add(201), []byte("Hello, Zinx!!!\n"))
	if err != nil {
		fmt.Println("Handle err:", err)
		return
	}
}

// StartHookFunc 给客户端发送hook执行成功的信息
func StartHookFunc(conn ziface.IConnection) {
	fmt.Println("StartHookFunc INVOKED")
	err := conn.SendMsg(202, []byte("StartHookFunc BEGIN\n"))
	if err != nil {
		fmt.Println(err)
		return
	}
	conn.SetAttribute("Author", "LewyHua")
	conn.SetAttribute("Github", "https://github.com/LewyHua")
}

// StopHookFunc 告诉其他的玩家这个用户下线了
func StopHookFunc(conn ziface.IConnection) {
	fmt.Println("StopHookFunc INVOKED")
	fmt.Printf("ConnID: %d is disconnected\n", conn.GetConnID())
	author, _ := conn.GetAttribute("Author")
	fmt.Printf("Author: %v\n", author)
	github, _ := conn.GetAttribute("Github")
	fmt.Printf("Author: %v\n", github)
}
