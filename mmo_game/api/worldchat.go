package api

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"zinx/mmo_game/core"
	"zinx/mmo_game/pb"
	"zinx/ziface"
	"zinx/znet"
)

type WorldChatApi struct {
	znet.BaseRouter
}

func (wc *WorldChatApi) Handle(req ziface.IRequest) {
	// 1 解析客户端传来的协议
	data := req.GetMsgData()
	protoMsg := &pb.Talk{}
	err := proto.Unmarshal(data, protoMsg)
	if err != nil {
		fmt.Println("Unmarshal Talk Msg err:", err)
		return
	}

	// 2 发送信息给所有当前世界的用户
	pid, err := req.GetConn().GetAttribute("PID")
	if err != nil {
		return
	}

	// 3 根据pid获取player对象
	id := pid.(int32)
	player := core.WorldMgrObj.GetPlayerByPID(id)

	// 4 广播
	content := protoMsg.GetContent()
	player.Talk(content)
}
