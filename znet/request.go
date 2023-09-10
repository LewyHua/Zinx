package znet

import (
	"zinx/ziface"
)

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好连接的conn
	data []byte             // 客户端请求的数据
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}

func NewRequest(conn ziface.IConnection, data []byte) ziface.IRequest {
	return &Request{
		conn: conn,
		data: data,
	}
}
