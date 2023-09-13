package ziface

import "net"

// IConnection 定义连接模块接口
type IConnection interface {
	Start()                                       // 启动连接 让当前连接准备开始工作
	Stop()                                        // 停止连接 结束连接的工作
	GetTCPConn() *net.TCPConn                     // 获取当前连接绑定的conn
	GetConnID() uint32                            // 获取当前连接ID
	GetRemoteAddr() net.Addr                      // 获取远程客户端TCP状态 IP Port
	SendMsg(msgID uint32, data []byte) error      // 把数据打包，发送给当前连接的worker的channel
	SetAttribute(key string, value interface{})   // 给conn设置自定义属性
	GetAttribute(key string) (interface{}, error) // 查询conn的自定义属性
	DelAttribute(key string)                      // 删除conn里面的自定义属性
}

// HandleFunc 定义一个处理连接业务的方法
type HandleFunc func(*net.TCPConn, []byte, int) error
