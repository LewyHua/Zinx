package ziface

// IServer 定义一个服务器接口
type IServer interface {
	Start()                                           // 启动服务器
	Stop()                                            // 停止服务器
	Serve()                                           // 运行服务器
	AddRouter(msgID uint32, router IRouter)           // 给当前服务注册一个路由方法，供客户端的连接处理使用
	GetConnManager() IConnManager                     // 获取当前服务器的连接管理器
	RegisterOnConnStart(func(connection IConnection)) // 注册连接建立后的hook函数
	RegisterOnConnStop(func(connection IConnection))  // 注册连接关闭前的hook函数
	InvokeOnConnStart(connection IConnection)         // 调用连接建立后的hook函数
	InvokeOnConnStop(connection IConnection)          // 调用连接关闭前的hook函数
}
