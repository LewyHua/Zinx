package main

import (
	"fmt"
	"zinx/mmo_game/core"
	"zinx/ziface"
	"zinx/znet"
)

func main() {
	s := znet.NewServer()
	// hook
	s.RegisterOnConnStart(OnConnected)

	// router api

	// start server
	s.Serve()

}

func OnConnected(conn ziface.IConnection) {
	player := core.NewPlayer(conn)
	// 发送Msg为1的PID消息
	player.SyncPID()
	// 发送Msg为200的出生位置消息
	player.BroadCastStartPosition()
	fmt.Printf("=====> PlayerID: %d is arrived at (%v, %v) \n", player.PID, player.X, player.Z)
}
