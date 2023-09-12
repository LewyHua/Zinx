package ziface

// IMsgHandler 消息管理层接口
type IMsgHandler interface {
	HandleMsg(request IRequest) error             // 执行对应的消息处理方法
	AddRouter(msgID uint32, router IRouter) error // 为消息添加具体业务逻辑
}
