package ziface

type IMessage interface {
	GetID() uint32
	GetLen() uint32
	GetData() []byte

	SetID(uint32)
	SetLen(uint32)
	SetData([]byte)
}
