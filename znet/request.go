package znet

import (
	"zinx/ziface"
)

type Request struct {
	conn ziface.IConnection // 已经和客户端建立好连接的conn
	msg  ziface.IMessage    // 客户端请求的数据
}

func (r *Request) GetConn() ziface.IConnection {
	return r.conn
}

func (r *Request) GetMsgData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetID()
}

func NewRequest(conn ziface.IConnection, msg ziface.IMessage) ziface.IRequest {
	return &Request{
		conn: conn,
		msg:  msg,
	}
}
