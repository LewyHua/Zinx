package ziface

// IMsgHandler 消息管理层接口
type IMsgHandler interface {
	HandleMsg(request IRequest)             // 执行对应的消息处理方法
	AddRouter(msgID uint32, router IRouter) // 为消息添加具体业务逻辑
	StartWorkerPool()                       // 启动Worker工作池
	SendMsgToTaskQueue(request IRequest)    // 发送消息到消息队列
}
