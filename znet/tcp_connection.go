package znet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"zinx/utils"
	"zinx/ziface"
)

type TCPConnection struct {
	Server     ziface.IServer         // 当前连接属于的Server
	Conn       *net.TCPConn           // 当前连接socket
	ConnID     uint32                 // 连接ID
	isClosed   bool                   // 连接状态
	ExitChan   chan bool              // 告知当前连接已经停止的channel
	MsgChan    chan []byte            // 读写channel之间的通信channel
	MsgHandler ziface.IMsgHandler     // 该连接处理的方法Router
	Attributes map[string]interface{} // 连接属性集合
	attrLock   sync.RWMutex           // 属性集合读写锁
}

func NewTCPConnection(server ziface.IServer, conn *net.TCPConn, connID uint32, msgHandler ziface.IMsgHandler) ziface.IConnection {
	c := &TCPConnection{
		Server:     server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		ExitChan:   make(chan bool, 1),
		MsgChan:    make(chan []byte),
		Attributes: make(map[string]interface{}),
	}
	c.Server.GetConnManager().AddConn(c)
	return c
}

// StartReader 当前连接的读业务方法
func (c *TCPConnection) StartReader() {
	fmt.Println("[Reader Goroutine running...]")
	defer fmt.Printf("Reader is exiting... ConnID: %d, RemoteAddr: %s\n", c.ConnID, c.GetRemoteAddr())
	defer c.Stop()

	for {
		// 创建pack对象
		dp := NewDataPack()

		// 读取客户端Msg Head 8 bytes
		headMsg := make([]byte, dp.GetHeadLen())
		_, err := io.ReadFull(c.Conn, headMsg)
		if err != nil {
			fmt.Println("server read head msg err:", err)
			break
		}

		// 把header解包到Message结构体
		msg, err := dp.Unpack(headMsg)
		if err != nil {
			fmt.Println("unpack head msg err:", err)
			break
		}

		var data []byte
		if msg.GetLen() > 0 {
			data = make([]byte, msg.GetLen())
			_, err := io.ReadFull(c.Conn, data)
			if err != nil {
				fmt.Println("server read data err:", err)
				break
			}
			msg.SetData(data)
		}

		// 得到当前conn以及数据的Request
		req := NewRequest(c, msg)

		// 如果工作池没有创建
		if utils.GlobalConfig.WorkerPoolSize > 0 {
			c.MsgHandler.SendMsgToTaskQueue(req)
		} else {
			// 从路由中找到注册绑定的Conn对应的Router调用
			go c.MsgHandler.HandleMsg(req)
		}
	}
}

// StartWriter 给客户端协消息模块
func (c *TCPConnection) StartWriter() {
	fmt.Println("[Writer goroutine running...]")
	defer fmt.Printf("Writer is exiting... ConnID: %d, RemoteAddr: %s\n", c.ConnID, c.GetRemoteAddr())

	for {
		select {
		case data := <-c.MsgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err:", err)
				return
			}
		case <-c.ExitChan:
			return
		}

	}
}

func (c *TCPConnection) Start() {
	fmt.Printf("Connection starting... ConnID = %d\n", c.ConnID)
	// 启动从当前连接读数据的业务
	go c.StartReader()

	// 启动写业务
	go c.StartWriter()

	// 按照开发者传递的建立连接之后的OnConnStart执行hook函数
	c.Server.InvokeOnConnStart(c)
}

func (c *TCPConnection) Stop() {
	fmt.Printf("Connection stopping... ConnID = %d\n", c.ConnID)
	if c.isClosed {
		return
	}
	c.isClosed = true

	// 按照开发者传递的关闭连接之前的OnConnStop执行hook函数
	c.Server.InvokeOnConnStop(c)

	// 关闭socket连接
	err := c.Conn.Close()
	if err != nil {
		fmt.Println("Conn.Close err:", err)
		return
	}

	c.ExitChan <- true

	// 从连接管理器删除当前conn
	c.Server.GetConnManager().DelConn(c)

	// 关闭channel
	close(c.ExitChan)
	close(c.MsgChan)
}

func (c *TCPConnection) GetConn() net.Conn {
	return c.Conn
}

func (c *TCPConnection) GetConnID() uint32 {
	return c.ConnID
}

func (c *TCPConnection) GetRemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *TCPConnection) SendMsg(msgID uint32, data []byte) error {
	if c.isClosed {
		return errors.New(fmt.Sprintf("Connection %d closed, cannot send data\n", msgID))
	}

	dp := NewDataPack()
	msgBinaries, err := dp.Pack(NewMessage(msgID, data))
	if err != nil {
		fmt.Printf("Server pack msg: %d err: %v\n", msgID, err)
		return errors.New("pack msg error")
	}

	c.MsgChan <- msgBinaries

	return nil
}

func (c *TCPConnection) SetAttribute(key string, value interface{}) {
	c.attrLock.Lock()
	defer c.attrLock.Unlock()
	c.Attributes[key] = value
}

func (c *TCPConnection) GetAttribute(key string) (value interface{}, err error) {
	c.attrLock.RLock()
	defer c.attrLock.RUnlock()
	value, ok := c.Attributes[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Attributes with key: %s not found", key))
	}
	return value, nil
}

func (c *TCPConnection) DelAttribute(key string) {
	c.attrLock.Lock()
	defer c.attrLock.Unlock()
	delete(c.Attributes, key)
}
