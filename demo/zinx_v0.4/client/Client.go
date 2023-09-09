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
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("Dial err:", err)
		return
	}

	// 2 写数据
	for {
		_, err = conn.Write([]byte("Hello zinx0.4"))
		if err != nil {
			fmt.Println("Write err:", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Read buf err:", err)
			return
		}
		fmt.Printf("Server call back: %s, len: %d\n", buf, cnt)
		time.Sleep(time.Second * 3)
	}
}
