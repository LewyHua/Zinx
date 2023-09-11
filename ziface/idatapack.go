package ziface

/*
	4B   4B
[DataLen|ID][Data]
\___Head___/\Body/
*/

// IDataPack 封包拆包接口
type IDataPack interface {
	GetHeadLen() uint32
	Pack(msg IMessage) ([]byte, error)
	Unpack([]byte) (IMessage, error)
}
