package core

import (
	"fmt"
	"github.com/gogo/protobuf/proto"
	"math/rand"
	"sync/atomic"
	"time"
	"zinx/mmo_game/pb"
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

// PID PIDGen ID Generator
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

// Talk Broadcast player chat
// 广播玩家聊天
func (p *Player) Talk(content string) {
	// 1. Assemble MsgID200 proto data
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  1, // TP: 1 represents chat broadcast (代表聊天广播)
		Data: &pb.BroadCast_Content{
			Content: content,
		},
	}

	// 2. Get all online players in the current world
	// 得到当前世界所有的在线玩家
	players := WorldMgrObj.GetAllPlayers()

	// 3. Send MsgID:200 message to all players
	// 向所有的玩家发送MsgID:200消息
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

func (p *Player) SyncSurrounding() {
	//1 根据自己的位置，获取周围九宫格内的玩家pID
	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	// 2 根据pID得到所有玩家对象
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
	}

	// 3 给这些玩家发送MsgID:200消息，让自己出现在对方视野中
	// 3.1 组建MsgID200 proto数据
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2, //TP:2 represents broadcasting coordinates (广播坐标)
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// 3.2 每个玩家分别给对应的客户端发送200消息，显示人物
	for _, player := range players {
		player.SendMsg(200, msg)
	}

	// 4 让周围九宫格内的玩家出现在自己的视野中
	// 4.1 制作Message SyncPlayers 数据
	playersData := make([]*pb.Player, 0, len(players))
	for _, player := range players {
		p := &pb.Player{
			PID: player.PID,
			P: &pb.Position{
				X: player.X,
				Y: player.Y,
				Z: player.Z,
				V: player.V,
			},
		}
		playersData = append(playersData, p)
	}

	// 4.2 封装SyncPlayer protobuf数据
	SyncPlayersMsg := &pb.SyncPlayers{
		Ps: playersData[:],
	}

	// 4.3 给当前玩家发送需要显示周围的全部玩家数据
	p.SendMsg(202, SyncPlayersMsg)
}

// UpdatePos Broadcast player position update
// (广播玩家位置移动)
func (p *Player) UpdatePos(x float32, y float32, z float32, v float32) {

	// 触发消失视野和添加视野业务
	// 计算旧格子gID
	oldGID := WorldMgrObj.AoiMgr.GetGIDByPos(p.X, p.Z)
	// 计算新格子gID
	newGID := WorldMgrObj.AoiMgr.GetGIDByPos(x, z)

	// 更新玩家的位置信息
	p.X = x
	p.Y = y
	p.Z = z
	p.V = v

	if oldGID != newGID {
		// 触发grid切换
		// 把pID从就的aoi格子中删除
		WorldMgrObj.AoiMgr.DelPIDFromGrid(int(p.PID), oldGID)

		// 把pID添加到新的aoi格子中去
		WorldMgrObj.AoiMgr.AddPIDToGrid(int(p.PID), newGID)

		_ = p.OnExchangeAoiGrid(oldGID, newGID)
	}

	// 组装protobuf协议，发送位置给周围玩家
	msg := &pb.BroadCast{
		PID: p.PID,
		Tp:  4, //Tp:4  Coordinates information after movement(移动之后的坐标信息)
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// (获取当前玩家周边全部玩家)
	players := p.GetSurroundingPlayers()

	// (向周边的每个玩家发送MsgID:200消息，移动位置更新消息)
	for _, player := range players {
		player.SendMsg(200, msg)
	}
}

func (p *Player) OnExchangeAoiGrid(oldGID, newGID int) error {
	// (获取就的九宫格成员)
	oldGrIDs := WorldMgrObj.AoiMgr.GetSurroundingGridsByGID(oldGID)

	// 为旧的九宫格成员建立哈希表,用来快速查找
	oldGrIDsMap := make(map[int]bool, len(oldGrIDs))
	for _, grID := range oldGrIDs {
		oldGrIDsMap[grID.GID] = true
	}

	// 获取新的九宫格成员
	newGrids := WorldMgrObj.AoiMgr.GetSurroundingGridsByGID(newGID)

	// 为新的九宫格成员建立哈希表,用来快速查找
	newGridsMap := make(map[int]bool, len(newGrids))
	for _, grID := range newGrids {
		newGridsMap[grID.GID] = true
	}

	//------ > Handle visibility disappearance (处理视野消失) <-------
	offlineMsg := &pb.SyncPID{
		PID: p.PID,
	}

	// (找到在旧的九宫格中出现,但是在新的九宫格中没有出现的格子)
	leavingGrIDs := make([]*Grid, 0)
	for _, grID := range oldGrIDs {
		if _, ok := newGridsMap[grID.GID]; !ok {
			leavingGrIDs = append(leavingGrIDs, grID)
		}
	}

	// (获取需要消失的格子中的全部玩家)
	for _, grID := range leavingGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)
		for _, player := range players {

			// Make oneself disappear in the views of other players
			// 让自己在其他玩家的客户端中消失
			player.SendMsg(201, offlineMsg)

			// Make other players' information disappear in one's own client
			// 将其他玩家信息 在自己的客户端中消失
			anotherOfflineMsg := &pb.SyncPID{
				PID: player.PID,
			}
			p.SendMsg(201, anotherOfflineMsg)
			time.Sleep(200 * time.Millisecond)
		}
	}

	// ------ > Handle visibility appearance(处理视野出现) <-------

	// Find the grid IDs that appear in the new nine-grid but not in the old nine-grid
	// 找到在新的九宫格内出现,但是没有在就的九宫格内出现的格子
	enteringGrIDs := make([]*Grid, 0)
	for _, grid := range newGrids {
		if _, ok := oldGrIDsMap[grid.GID]; !ok {
			enteringGrIDs = append(enteringGrIDs, grid)
		}
	}

	onlineMsg := &pb.BroadCast{
		PID: p.PID,
		Tp:  2,
		Data: &pb.BroadCast_P{
			P: &pb.Position{
				X: p.X,
				Y: p.Y,
				Z: p.Z,
				V: p.V,
			},
		},
	}

	// Get all players in the appearing grids
	// 获取需要显示格子的全部玩家
	for _, grID := range enteringGrIDs {
		players := WorldMgrObj.GetPlayersByGID(grID.GID)

		for _, player := range players {
			// Make oneself appear in the views of other players
			// 让自己出现在其他人视野中
			player.SendMsg(200, onlineMsg)

			// Make other players appear in one's own client
			// 让其他人出现在自己的视野中
			anotherOnlineMsg := &pb.BroadCast{
				PID: player.PID,
				Tp:  2,
				Data: &pb.BroadCast_P{
					P: &pb.Position{
						X: player.X,
						Y: player.Y,
						Z: player.Z,
						V: player.V,
					},
				},
			}

			time.Sleep(200 * time.Millisecond)
			p.SendMsg(200, anotherOnlineMsg)
		}
	}

	return nil
}

// GetSurroundingPlayers Get information of surrounding players in the current player's AOI
// 获得当前玩家的AOI周边玩家信息
func (p *Player) GetSurroundingPlayers() []*Player {
	// Get all pIDs in the current AOI area
	// 得到当前AOI区域的所有pID
	pIDs := WorldMgrObj.AoiMgr.GetPIDsByPos(p.X, p.Z)

	// Put all players corresponding to pIDs into the Player slice
	// 将所有pID对应的Player放到Player切片中
	players := make([]*Player, 0, len(pIDs))
	for _, pID := range pIDs {
		players = append(players, WorldMgrObj.GetPlayerByPID(int32(pID)))
	}

	return players
}

// LostConnection logs off
// 玩家下线
func (p *Player) LostConnection() {
	// 1 Get players in the surrounding AOI nine-grid
	// 获取周围AOI九宫格内的玩家
	players := p.GetSurroundingPlayers()

	// 2 Assemble MsgID:201 message
	// 封装MsgID:201消息
	msg := &pb.SyncPID{
		PID: p.PID,
	}

	// 3 Send messages to surrounding players
	// 向周围玩家发送消息
	for _, player := range players {
		player.SendMsg(201, msg)
	}

	// 4 Remove the current player from AOI in the world manager
	// 世界管理器将当前玩家从AOI中摘除
	WorldMgrObj.AoiMgr.DelFromGridByPos(int(p.PID), p.X, p.Z)
	WorldMgrObj.RemovePlayerByPID(p.PID)
}
