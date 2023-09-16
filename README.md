# Zinx

Zinx是一款轻量级Golang实现的游戏服务器框架，支持TCP/KCP连接，可自定义连接建立与关闭的hook函数，可自定义Router处理函数。
可通过配置来按需设置服务器工作池大小。\
详情参考：https://github.com/aceld 同名项目

### Demo 如何启动？
```shell
# 1. 进入demo目录 
cd demo 
# 2. 进入config目录，按需修改config.yaml的配置
vim config/config.yaml
# 3. 启动Server服务
go run zinx_v0.10/server/Server.go
# 4. 启动Client服务
go run zinx_v0.10/client1/Client1.go
```

### 如何使用？
```go
func main() {
	// 1 创建一个server句柄，使用zinx的api
	s := znet.NewServer()

	// 2 注册开启关闭conn时的hook函数
	//   需要时自行创建即可函数签名为func(conn ziface.IConnection)的方法即可
	s.RegisterOnConnStart(StartHookFunc)
	s.RegisterOnConnStop(StopHookFunc)

	// 3 给当前Zinx框架增加自定义Router
	//   Router是任何实现了IRouter接口的实现类
	//   提供了模版方法模式，通过继承BaseRouter的任意方法即可，例如下面的PingRouter
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloRouter{})
	
	// 4 启动服务
	s.Serve()
}

// PingRouter 自定义路由
type PingRouter struct {
    znet.BaseRouter
}

func (p *PingRouter) Handle(request ziface.IRequest) {
   var atomicInt atomic.Uint32
   fmt.Println("Calling PingRouter...")
   fmt.Printf("MsgID: %d, Data: %v\n", request.GetMsgID(), string(request.GetMsgData()))
   err := request.GetConn().SendMsg(atomicInt.Add(200), []byte("Pong...\n"))
   if err != nil {
	   fmt.Println("Handle err:", err)
   return
   }
}
```

### 应用服务流程
1. Server启动服务
   1. 根据开发者创建server实例时指定的路由方法注册路由/hook函数
   2. 通过MsgHandler创建channel消息队列，以及指定数量的workerPool等待处理消息
   3. 创建配置文件指定/默认工作池
   4. 根据配置文件开启对应TCP/KCP监听端口，等待客户端连接
2. 客户端启动服务
   1. 客户端通过DataPack组件、或自行打包消息
   2. 客户端发送打包好后的二进制字节流到服务端
   3. 客户端监听等待接受服务端消息
3. Server接收到客户端链接
   1. 封装*net.Conn称为TCPConnection/KCPConnection
   2. 注入Server、连接ID、Router到Connection
4. 启动Connection开始业务
   1. Connection开启Reader协程
      1. Reader方法读取客户端输入
      2. 解决粘包，读取前8字节获取消息Header
      3. 读取消息体，一起封装到Message结构体
      4. 封装Message和Connection成为一个Request
      5. 通过MsgHandler将消息发送到消息队列里面
   2. Connection开启Writer协程
      1. 阻塞从消息队列等待获取消息，收到消息发送给客户端
      2. 如果收到退出信号，则关闭Writer协程
5. MsgHandler接收到Request，通过Request获取连接ID取模获取workerID，并发送到对应worker接收的消息队列
6. 传入Request作为参数，worker接收到request解析出MsgId对应的路由方法，调用Connection.Router的三个hook方法
7. router

### Server 服务器
```go
// IServer 定义一个服务器接口
type IServer interface {
    Start()                                           // 启动服务器
    Stop()                                            // 停止服务器
    Serve()                                           // 运行服务器
    AddRouter(msgID uint32, router IRouter)           // 给当前服务注册一个路由方法，供客户端的连接处理使用
    GetConnManager() IConnManager                     // 获取当前服务器的连接管理器
    RegisterOnConnStart(func(connection IConnection)) // 注册连接建立后的hook函数
    RegisterOnConnStop(func(connection IConnection))  // 注册连接关闭前的hook函数
    InvokeOnConnStart(connection IConnection)         // 调用连接建立后的hook函数
    InvokeOnConnStop(connection IConnection)          // 调用连接关闭前的hook函数
}
```

### Connection TCP/KCP连接
```go
// IConnection 定义连接模块接口
type IConnection interface {
    Start()                                       // 启动连接 让当前连接准备开始工作
    Stop()                                        // 停止连接 结束连接的工作
    GetConn() net.Conn                            // 获取当前连接绑定的conn
    GetConnID() uint32                            // 获取当前连接ID
    GetRemoteAddr() net.Addr                      // 获取远程客户端TCP状态 IP Port
    SendMsg(msgID uint32, data []byte) error      // 把数据打包，发送给当前连接的worker的channel
    SetAttribute(key string, value interface{})   // 给conn设置自定义属性
    GetAttribute(key string) (interface{}, error) // 查询conn的自定义属性
    DelAttribute(key string)                      // 删除conn里面的自定义属性
}
```

### Connection Manager 连接管理模块
```go
type IConnManager interface {
    AddConn(conn IConnection)
    DelConn(conn IConnection) // TODO DelConnByID
    GetConnByID(connID uint32) (IConnection, error)
    GetConnNum() int
    ClearConns()
}
```

### Request 封装请求
```go
type Request struct {
    conn ziface.IConnection // 已经和客户端建立好连接的conn
    data []byte             // 客户端请求的数据
}
```

### Message 消息模块
```go
type Request struct {
    conn ziface.IConnection // 已经和客户端建立好连接的conn
    data []byte             // 客户端请求的数据
}
```

### Message Handler 消息管理模块
```go
type IMsgHandler interface {
    HandleMsg(request IRequest)             // 执行对应的消息处理方法
    AddRouter(msgID uint32, router IRouter) // 为消息添加具体业务逻辑
    StartWorkerPool()                       // 启动Worker工作池
    SendMsgToTaskQueue(request IRequest)    // 发送消息到消息队列
}
```

### Router 定义路由方法
```go
type BaseRouter struct {
}
func (b *BaseRouter) PreHandle(request ziface.IRequest)  {}
func (b *BaseRouter) Handle(request ziface.IRequest)     {}
func (b *BaseRouter) PostHandle(request ziface.IRequest) {}

```

### Data Pack 数据封装模块
```go
// IDataPack 封包拆包接口
type IDataPack interface {
    GetHeadLen() uint32
    Pack(msg IMessage) ([]byte, error)
    Unpack([]byte) (IMessage, error)
}

```