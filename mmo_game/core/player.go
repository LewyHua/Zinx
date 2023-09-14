package core

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"sync/atomic"
	"zinx/mmo_game/pb/pb"
	"zinx/ziface"
)

type Player struct {
	PID  int32              // 玩家ID
	Conn ziface.IConnection // 玩家的连接
	X    float32            // Planar x coordinate(平面x坐标)
	Y    float32            // Height(高度)
	Z    float32            // Planar y coordinate (Note: not Y)- 平面y坐标 (注意不是Y)
	V    float32            // Rotation 0-360 degrees(旋转0-360度)
}

// PIDGen ID Generator
// var PIDGen uint32 = 1 // Counter for generating player IDs(用来生成玩家ID的计数器)
// var IDLock sync.Mutex // Mutex for protecting PIDGen(保护PIDGen的互斥机制)
var PID atomic.Int32

// NewPlayer Create a player object
func NewPlayer(conn ziface.IConnection) *Player {
	ID := PID.Add(1)
	p := &Player{
		PID:  ID,
		Conn: conn,
		X:    float32(160 + rand.Intn(50)), // Randomly offset on the X-axis based on the point 160(随机在160坐标点 基于X轴偏移若干坐标)
		Y:    0,                            // Height is 0
		Z:    float32(134 + rand.Intn(50)), // Randomly offset on the Y-axis based on the point 134(随机在134坐标点 基于Y轴偏移若干坐标)
		V:    0,                            // Angle is 0, not yet implemented(角度为0，尚未实现)
	}
	return p
}

// SyncPID Inform the client about pID and synchronize the generated player ID to the client
// (告知客户端pID,同步已经生成的玩家ID给客户端)
func (p *Player) SyncPID() {
	// Assemble MsgID0 proto data
	// (组建MsgID0 proto数据)
	data := &pb.SyncPID{
		PID: p.PID,
	}

	// Send data to the client
	// (发送数据给客户端)
	p.SendMsg(1, data)
}

// BroadCastStartPosition Broadcast the player's starting position
// (广播玩家自己的出生地点)
func (p *Player) BroadCastStartPosition() {

	// Assemble MsgID200 proto data
	// (组建MsgID200 proto数据)
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2, //TP:2  represents broadcasting coordinates (广播坐标)
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// Send data to the client
	// 发送数据给客户端
	p.SendMsg(200, msg)
}

// SendMsg Send messages to the client, mainly serializing and sending the protobuf data of the pb Message
//
//	(发送消息给客户端，主要是将pb的protobuf数据序列化之后发送)
func (p *Player) SendMsg(msgID uint32, data proto.Message) {
	if p.Conn == nil {
		fmt.Println("connection in player is nil")
		return
	}

	// fmt.Printf("before Marshal data = %+v\n", data)

	// Serialize the proto Message structure
	// 将proto Message结构体序列化
	msg, err := proto.Marshal(data)
	if err != nil {
		fmt.Println("marshal msg err: ", err)
		return
	}

	// fmt.Printf("after Marshal data = %+v\n", msg)

	// Call the Zinx framework's SendMsg to send the packet
	// 调用Zinx框架的SendMsg发包
	if err := p.Conn.SendMsg(msgID, msg); err != nil {
		fmt.Println("Player SendMsg error !")
		return
	}

	return
}
