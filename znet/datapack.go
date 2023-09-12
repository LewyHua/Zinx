package znet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type DataPack struct{}

func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

// GetHeadLen 包Header长度
func (dp *DataPack) GetHeadLen() uint32 {
	// Len: uint32(4byte)  ID: uint32(4byte)
	return 8
}

// Pack 封包 |DataLen|ID|Data|
func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建存放bytes的字节缓冲区
	buf := bytes.NewBuffer([]byte{})

	// DataLen写入buf
	err := binary.Write(buf, binary.LittleEndian, msg.GetLen())
	if err != nil {
		return nil, err
	}

	// ID写入buf
	err = binary.Write(buf, binary.LittleEndian, msg.GetID())
	if err != nil {
		return nil, err
	}

	// Data写入buf
	err = binary.Write(buf, binary.LittleEndian, msg.GetData())
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// Unpack 拆包 读出Head
func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建从输入读取二进制数据的ioReader
	buf := bytes.NewReader(binaryData)
	// 解包Head 得到len和id
	msg := new(Message)

	// 读取Len
	err := binary.Read(buf, binary.LittleEndian, &msg.Len)
	if err != nil {
		return nil, err
	}

	// 读取ID
	err = binary.Read(buf, binary.LittleEndian, &msg.ID)
	if err != nil {
		return nil, err
	}

	// 判断Len是否超出package最大允许长度
	if msg.Len > utils.GlobalConfig.MaxPackageSize {
		return nil, errors.New(fmt.Sprintf("Package size too large, max: %v", utils.GlobalConfig.MaxConn))
	}

	return msg, nil
}
