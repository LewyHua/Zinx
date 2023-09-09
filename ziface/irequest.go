package ziface

// IRequest 把客户端请求连接和数据包装成一个请求
type IRequest interface {
	GetConn() IConnection // 获取当前连接
	GetData() []byte      // 获取请求的消息数据
}
