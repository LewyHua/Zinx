# Zinx

### 如何启动？
```shell
# 1. 进入demo目录 
cd demo 
# 2. 进入config目录，按需修改config.yaml的配置
vim config/config.yaml
# 3. 启动Server服务
go run server/Server.go
# 4. 启动Client服务
go run client/Client.go
```

### 应用服务流程
1. Server启动服务，监听端口，等待客户端连接
2. 客户端请求进入Server 
3. Server封装Connection，注入Router到Connection
4. 启动Connection开始业务
5. Connection开启Reader协程
6. Reader方法读取客户端输入，封装数据和Connection成为一个Request
7. 传入Request作为参数，调用Connection.Router的三个hook方法

### Server 服务器
```go
type Server struct {
    Name      string         // 服务器名称
    IPVersion string         // 服务器绑定IP版本
    IP        string         // 服务器绑定的IP
    Port      int            // 服务器监听端口
    Router    ziface.IRouter // 当前server连接注册的对应处理业务
}
```

### Connection TCP连接
```go
type Connection struct {
    Conn     *net.TCPConn   // 当前连接socket
    ConnID   uint32         // 连接ID
    isClosed bool           // 连接状态
    ExitChan chan bool      // 告知当前连接已经停止的channel
    Router   ziface.IRouter // 该连接处理的方法Router
}
```

### Request 封装请求
```go
type Request struct {
    conn ziface.IConnection // 已经和客户端建立好连接的conn
    data []byte             // 客户端请求的数据
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