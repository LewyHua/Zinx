package main

import (
	"fmt"
	"zinx/mmo_game/api"
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	s := znet.NewServer()
	// hook
	s.RegisterOnConnStart(OnConnected)
	s.RegisterOnConnStop(OnDisconnected)

	// router api
	s.AddRouter(2, &api.WorldChatApi{})
	s.AddRouter(3, &api.MoveApi{})

	// start server
	s.Serve()

}

func OnConnected(conn ziface.IConnection) {
	player := core.NewPlayer(conn)
	// 给当前连接绑定PID
	conn.SetAttribute("PID", player.PID)
	// 发送Msg为1的PID消息
	player.SyncPID()
	// 发送Msg为200的出生位置消息
	player.BroadCastStartPosition()
	// 广播当前用户上线的消息到世界
	core.WorldMgrObj.AddPlayer(player)
	// 告知周边玩家当前玩家已经上线（广播当前玩家位置信息）
	player.SyncSurrounding()
	// 告知当前玩家周边玩家的信息

	fmt.Printf("=====> PlayerID: %d is arrived at (%v, %v) \n", player.PID, player.X, player.Z)
}

func OnDisconnected(conn ziface.IConnection) {
	PID, err := conn.GetAttribute("PID")
	if err != nil {
		conn.Stop()
		return
	}
	player := core.WorldMgrObj.GetPlayerByPID(PID.(int32))
	player.LostConnection()
	fmt.Printf("=====> PlayerID: %d disconnected \n", player.PID)

}
