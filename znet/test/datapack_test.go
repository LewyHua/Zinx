package test

import (
	"fmt"
	"io"
	"net"
	"testing"
	"zinx/znet"
)

func TestDataPack(t *testing.T) {
	/*
		模拟服务器
	*/
	// 1 创建socket
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err:", err)
	}

	// 2 从客户端读取，拆包 协程处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept err:", err)
				//return
			}
			go func(conn net.Conn) {
				// 处理客户端请求
				// 拆包
				dp := znet.NewDataPack()
				for {
					// 存放head数据
					headData := make([]byte, dp.GetHeadLen())
					// 获取head数据
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head err:", err)
						break
					}

					// head数据解包到msg
					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("server unpack head err:", err)
						return
					}

					// 有数据
					if msgHead.GetLen() > 0 {
						msg := msgHead.(*znet.Message)
						msg.Data = make([]byte, msg.GetLen())

						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err:", err)
							return
						}
						fmt.Printf("Received message: %#v\n", msg)
					}

				}

			}(conn)
		}
	}()

	/*
		模拟客户端
	*/
	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client1 dial err", err)
		return
	}

	// 发送两个msg
	dp := znet.NewDataPack()
	msg1 := &znet.Message{
		ID:   1,
		Len:  5,
		Data: []byte("Hello"),
	}
	bytes1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("pack msg1 err:", err)
		return
	}

	msg2 := &znet.Message{
		ID:   2,
		Len:  4,
		Data: []byte("Zinx"),
	}
	bytes2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 err:", err)
		return
	}

	bytes1 = append(bytes1, bytes2...)
	_, err = conn.Write(bytes1)
	if err != nil {
		fmt.Println("client1 send err:", err)
		return
	}

	select {}
}
