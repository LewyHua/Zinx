package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("Client starting...")
	// 1 链接服务端，获取conn
	conn, err := net.Dial("udp", "255.255.255.255:12345")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}

	// 2 写数据
	for {
		//// 客户端封包
		//dp := znet.NewDataPack()
		//bytes, err := dp.Pack(znet.NewMessage(1, []byte("Client Hello!!!")))
		//if err != nil {
		//	fmt.Println("client1 pack msg err:", err)
		//	return
		//}
		// 客户端发数据
		if _, err = conn.Write([]byte("asdaasdadasdadsadasdas")); err != nil {
			fmt.Println("client1 write err:", err)
			return
		}
		//fmt.Println("sent")

		// 客户端接收数据，解包头
		//binaryHead := make([]byte, 8)
		//if _, err = io.ReadFull(conn, binaryHead); err != nil {
		//	fmt.Println("Read Head err:", err)
		//}
		//
		//message, err := dp.Unpack(binaryHead)
		//if err != nil {
		//	fmt.Println("Unpack head err:", err)
		//	return
		//}

		//if message.GetLen() > 0 {
		//	binaryBody := make([]byte, message.GetLen())
		//	if _, err := io.ReadFull(conn, binaryBody); err != nil {
		//		fmt.Println("Unpack body err:", err)
		//		return
		//	}
		//	message.SetData(binaryBody)
		//	fmt.Printf("Received from Server ---> ID: %d, Len: %d, Data: %v", message.GetID(), message.GetLen(), string(message.GetData()))
		//}

		time.Sleep(time.Millisecond)
	}
}
