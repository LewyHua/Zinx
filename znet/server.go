package znet

import (
	"errors"
	"fmt"
	"net"
	"sync/atomic"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name      string // 服务器名称
	IPVersion string // 服务器绑定IP版本
	IP        string // 服务器绑定的IP
	Port      int    // 服务器监听端口
}

// CallBackToClient 定义当前客户端连接所绑定的handle api
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[Conn Handle] Callback to client")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back err:", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start() {
	fmt.Printf("[Start] Server Listener at %s:%d is starting...\n", s.IP, s.Port)

	go func() {

		// 1 获取tcp的addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr err:", err)
			return
		}

		// 2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP err:", err)
			return
		}
		fmt.Printf("Started Zinx at %s:%d success, Listening...\n", s.IP, s.Port)

		var cid atomic.Uint32

		// 3 阻塞地等待客户端链接，处理客户端链接业务
		for {
			// 3.1 有客户端链接进来，阻塞会返回
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP err:", err)
				return
			}

			dealConn := NewConnection(conn, cid.Add(1), CallBackToClient)
			go dealConn.Start()
		}

	}()
}

func (s *Server) Stop() {
	// TODO 对服务器资源回收/停止
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 启动server后的额外业务

	// 阻塞
	select {}
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
	return s
}
