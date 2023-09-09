package znet

import (
	"fmt"
	"io"
	"net"
	"zinx/ziface"
)

type Connection struct {
	Conn      *net.TCPConn      // 当前连接socket
	ConnID    uint32            // 连接ID
	isClosed  bool              // 连接状态
	HandleAPI ziface.HandleFunc // 连接绑定的业务方法API
	ExitChan  chan bool         // 告知当前连接已经停止的channel
}

func NewConnection(conn *net.TCPConn, connID uint32, callbackApi ziface.HandleFunc) ziface.IConnection {
	return &Connection{
		Conn:      conn,
		ConnID:    connID,
		isClosed:  false,
		HandleAPI: callbackApi,
		ExitChan:  make(chan bool, 1),
	}
}

// StartReader 当前连接的读业务方法
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine running...")
	defer fmt.Printf("Reader is exiting... ConnID: %d, RemoteAddr: %s", c.ConnID, c.GetRemoteAddr())
	defer c.Stop()

	for {
		// 读取客户端数据到buf，最大512字节
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil && err == io.EOF {
			fmt.Printf("Client ConnID:%d closed\n", c.ConnID)
			return
		}
		if err != nil {
			fmt.Println("Receive msg err:", err)
			continue
		}

		// 调用当前连接所绑定的HandleAPI
		err = c.HandleAPI(c.Conn, buf, cnt)
		if err != nil {
			fmt.Printf("ConnID: %d, HandleAPI err: %s\n", c.ConnID, err)
			continue
		}
	}
}

func (c *Connection) Start() {
	fmt.Printf("Connection starting... ConnID = %d\n", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()

	// TODO 启动写业务
}

func (c *Connection) Stop() {
	fmt.Printf("Connection stopping... ConnID = %d\n", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn.Close err:", err)
		return
	}

	// 关闭channel
	close(c.ExitChan)
}

func (c *Connection) GetTCPConn() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) Send(data []byte) bool {
	return false
}
