package znet

type Message struct {
	ID   uint32 // 消息ID
	Len  uint32 // 消息长度
	Data []byte // 消息内容
}
