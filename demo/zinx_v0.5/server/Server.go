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
	s := znet.NewServer("[zinx v0.5]")

	// 2 给当前Zinx框架增加自定义Router
	s.AddRouter(&PingRouter{})
	s.Serve()
}

// PingRouter test自定义路由
type PingRouter struct {
	znet.BaseRouter
}

//func (p *PingRouter) PreHandle(request ziface.IRequest) {
//	fmt.Println("PreHandling...")
//	_, err := request.GetConn().GetTCPConn().Write([]byte("Before Ping...\n"))
//	if err != nil {
//		fmt.Println("PreHandle err:", err)
//		return
//	}
//}

func (p *PingRouter) Handle(request ziface.IRequest) {
	var atomicInt atomic.Uint32
	fmt.Println("Handling...")
	fmt.Printf("MsgID: %d, Data: %v", request.GetMsgID(), string(request.GetMsgData()))
	err := request.GetConn().SendMsg(atomicInt.Add(1), []byte("Ping...\n"))
	if err != nil {
		fmt.Println("Handle err:", err)
		return
	}
}

//func (p *PingRouter) PostHandle(request ziface.IRequest) {
//	fmt.Println("PostHandling...")
//	_, err := request.GetConn().GetTCPConn().Write([]byte("After Ping...\n"))
//	if err != nil {
//		fmt.Println("PostHandle err:", err)
//		return
//	}
//}
