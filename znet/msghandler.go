package znet

import (
	"fmt"
	"zinx/utils"
	"zinx/ziface"
)

type MsgHandler struct {
	Apis           map[uint32]ziface.IRouter // 路由方法表
	TaskQueue      []chan ziface.IRequest    // 负责Worker取任务的消息队列
	WorkerPoolSize uint32                    // 业务工作Worker池的worker数量

}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           map[uint32]ziface.IRouter{},
		WorkerPoolSize: utils.GlobalConfig.WorkerPoolSize,
		TaskQueue:      make([]chan ziface.IRequest, utils.GlobalConfig.WorkerPoolSize),
	}
}

func (mh *MsgHandler) HandleMsg(request ziface.IRequest) {
	// 获取router
	router, ok := mh.Apis[request.GetMsgID()]
	// router不存在
	if !ok {
		fmt.Printf("router with msgID: %d not exists", request.GetMsgID())
		return
	}
	// 调用router
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		fmt.Printf("router with msgID: %d updated", msgID)
	}
	mh.Apis[msgID] = router
	fmt.Printf("Added a new router with msgID: %d success\n", msgID)
}

// StartWorkerPool 启动Worker工作池，一个Zinx只有一个工作池
func (mh *MsgHandler) StartWorkerPool() {
	// 根据WorkerPoolSize开启Worker，每个Worker通过协程承载
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// 启动一个Worker
		// 给每个worker对应的channel（消息队列）开辟空间
		mh.TaskQueue[i] = make(chan ziface.IRequest, utils.GlobalConfig.MaxWorkerTaskLen)
		// 启动当前Worker，阻塞等待channel中的消息
		go mh.StartWorker(i, mh.TaskQueue[i])
	}
}

// StartWorker 启动一个Worker的工作流程
func (mh *MsgHandler) StartWorker(workerID int, taskQueue chan ziface.IRequest) {
	fmt.Printf("Work ID: %v is starting...\n", workerID)
	for {
		select {
		// 如果有消息，就是客户端的一个request，执行其绑定的业务
		case request := <-taskQueue:
			mh.HandleMsg(request)
		}
	}

}

// SendMsgToTaskQueue 讲消息交给TaskQueue，Worker进行处理
func (mh *MsgHandler) SendMsgToTaskQueue(request ziface.IRequest) {
	// 1 将消息平均分配到不同的Worker
	workerId := request.GetConn().GetConnID() % mh.WorkerPoolSize
	fmt.Printf("Added ConnID: %v request to WorkerID: %v\n", request.GetConn().GetConnID(), workerId)

	// 2 将消息发送给worker对应的queue
	mh.TaskQueue[workerId] <- request
}
