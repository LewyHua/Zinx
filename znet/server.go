package znet

import (
	"fmt"
	"github.com/xtaci/kcp-go"
	"net"
	"sync/atomic"
	"zinx/utils"
	"zinx/ziface"
)

// Server IServer的接口实现，定义一个Server的服务器模块
type Server struct {
	Name        string                              // 服务器名称
	IPVersion   string                              // 服务器绑定IP版本
	IP          string                              // 服务器绑定的IP
	Port        int                                 // 服务器TCP监听端口
	MsgHandler  ziface.IMsgHandler                  // 当前server连接注册的对应处理业务
	ConnManager ziface.IConnManager                 // 当前server的连接管理器
	OnConnStart func(connection ziface.IConnection) //创建连接后的hook方法
	OnConnStop  func(connection ziface.IConnection) //关闭连接前的hook方法
	CID         atomic.Uint32                       // 全局连接id
}

func (s *Server) Start() {
	fmt.Printf("[Zinx] Server is starting...\n")
	fmt.Printf("[Zinx] Server Name: %s, Version: %s, Listen at %s:%d\n", s.Name, utils.GlobalConfig.Version, s.IP, s.Port)
	// 0 开启消息队列以及Worker工作池
	s.MsgHandler.StartWorkerPool()

	switch utils.GlobalConfig.Mode {
	case "tcp":
		go s.ListenTCPConn()
	case "kcp":
		go s.ListenKCPConn()
	default:
		go s.ListenTCPConn()
	}
}

func (s *Server) ListenTCPConn() {
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
	fmt.Printf("Started Zinx at %s:%d success, Listening using TCP ...\n", s.IP, s.Port)

	// 3 阻塞地等待客户端链接，处理客户端链接业务
	for {
		// 3.1 有客户端链接进来，阻塞会返回
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("AcceptTCP err:", err)
			return
		}

		// 判断是否超过最大连接数, 若超过，则关闭此次连接
		if s.ConnManager.GetConnNum() >= utils.GlobalConfig.MaxConn {
			// TODO 给客户端响应错误信息 / 让客户端等待
			fmt.Printf("Exceed max connection size: %d, connection stopped\n", utils.GlobalConfig.MaxConn)
			conn.Close()
			continue
		}

		dealConn := NewTCPConnection(s, conn, s.CID.Add(1), s.MsgHandler)
		go dealConn.Start()
	}
}

func (s *Server) ListenKCPConn() {
	// 1 监听KCP服务器地址
	listener, err := kcp.Listen(fmt.Sprintf("%s:%d", s.IP, s.Port))
	if err != nil {
		fmt.Println("Listen KCP err:", err)
		return
	}
	fmt.Printf("Started Zinx at %s:%d success, Listening using KCP...\n", s.IP, s.Port)

	// 2 阻塞地等待客户端链接，处理客户端链接业务
	for {
		// 2.1 有客户端链接进来，阻塞会返回
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("AcceptTCP err:", err)
			return
		}

		// 判断是否超过最大连接数, 若超过，则关闭此次连接
		if s.ConnManager.GetConnNum() >= utils.GlobalConfig.MaxConn {
			// TODO 给客户端响应错误信息 / 让客户端等待
			fmt.Printf("Exceed max connection size: %d, connection stopped\n", utils.GlobalConfig.MaxConn)
			conn.Close()
			continue
		}

		dealConn := NewKCPConnection(s, conn.(*kcp.UDPSession), s.CID.Add(1), s.MsgHandler)
		go dealConn.Start()
	}
}

func (s *Server) Stop() {
	// 对服务器资源回收/停止
	s.ConnManager.ClearConns()
}

func (s *Server) Serve() {
	// 启动server的服务功能
	s.Start()

	// TODO 启动server后的额外业务

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgID uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
}

func (s *Server) RegisterOnConnStart(onStartFunc func(connection ziface.IConnection)) {
	s.OnConnStart = onStartFunc
}

func (s *Server) RegisterOnConnStop(onStopFunc func(connection ziface.IConnection)) {
	s.OnConnStop = onStopFunc
}

func (s *Server) InvokeOnConnStart(connection ziface.IConnection) {
	if s.OnConnStart != nil {
		fmt.Println("Invoking OnConnStart()...")
		s.OnConnStart(connection)
	}
}

func (s *Server) InvokeOnConnStop(connection ziface.IConnection) {
	if s.OnConnStop != nil {
		fmt.Println("Invoking OnConnStop()...")
		s.OnConnStop(connection)
	}
}

func (s *Server) GetConnManager() ziface.IConnManager {
	return s.ConnManager
}

func NewServer(options ...Option) ziface.IServer {
	s := &Server{
		Name:        utils.GlobalConfig.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalConfig.Host,
		Port:        utils.GlobalConfig.Port,
		MsgHandler:  NewMsgHandler(),
		ConnManager: NewConnManager(),
	}

	for _, option := range options {
		option(s)
	}

	return s
}
