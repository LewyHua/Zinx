package znet

import (
	"errors"
	"fmt"
	"zinx/ziface"
)

type MsgHandler struct {
	Apis map[uint32]ziface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: map[uint32]ziface.IRouter{},
	}
}

func (mh *MsgHandler) HandleMsg(request ziface.IRequest) error {
	// 获取router
	router, ok := mh.Apis[request.GetMsgID()]
	// router不存在
	if !ok {
		return errors.New(fmt.Sprintf("router with msgID: %d not exists", request.GetMsgID()))

	}
	// 调用router
	router.PreHandle(request)
	router.Handle(request)
	router.PostHandle(request)
	return nil
}

func (mh *MsgHandler) AddRouter(msgID uint32, router ziface.IRouter) error {
	if _, ok := mh.Apis[msgID]; ok {
		return errors.New(fmt.Sprintf("router with msgID: %d already exists", msgID))
	}
	mh.Apis[msgID] = router
	fmt.Printf("Added a new router with msgID: %d success\n", msgID)
	return nil
}
