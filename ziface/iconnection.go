package ziface

import "net"

// IConnection 定义连接模块接口
type IConnection interface {
	Start()                   // 启动连接 让当前连接准备开始工作
	Stop()                    // 停止连接 结束连接的工作
	GetTCPConn() *net.TCPConn // 获取当前连接绑定的conn
	GetConnID() uint32        // 获取当前连接ID
	GetRemoteAddr() net.Addr  // 获取远程客户端TCP状态 IP Port
	Send(data []byte) bool    // 发送数据
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
