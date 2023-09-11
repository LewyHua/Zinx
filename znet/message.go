package znet

type Message struct {
	ID   uint32 // 消息ID
	Len  uint32 // 消息长度
	Data []byte // 消息内容
}

func NewMessage(id uint32, data []byte) *Message {
	return &Message{
		ID:   id,
		Len:  uint32(len(data)),
		Data: data,
	}
}

func (m *Message) GetID() uint32 {
	return m.ID
}

func (m *Message) GetLen() uint32 {
	return m.Len
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) SetID(id uint32) {
	m.ID = id
}

func (m *Message) SetLen(len uint32) {
	m.Len = len
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}
