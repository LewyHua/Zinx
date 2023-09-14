package core

import (
	"fmt"
	"sync"
)

// Grid 一个AOI地图中的格子类型
type Grid struct {
	GID       int          //格子ID
	MinX      int          //格子左边界
	MaxX      int          //格子右边界
	MinY      int          //格子上边界坐标
	MaxY      int          //格子下边界坐标
	playerIDs map[int]bool //当前格子内的玩家或者物体成员ID
	pIDLock   sync.RWMutex //playerIDs的保护map的锁
}

// NewGrid 初始化一个格子
func NewGrid(gID, minX, maxX, minY, maxY int) *Grid {
	return &Grid{
		GID:       gID,
		MinX:      minX,
		MaxX:      maxX,
		MinY:      minY,
		MaxY:      maxY,
		playerIDs: make(map[int]bool),
	}
}

// Add 向当前格子中添加一个玩家
func (g *Grid) Add(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	g.playerIDs[playerID] = true
}

// Del 从格子中删除一个玩家
func (g *Grid) Del(playerID int) {
	g.pIDLock.Lock()
	defer g.pIDLock.Unlock()

	delete(g.playerIDs, playerID)
}

// GetPlayerIDs 得到当前格子的所有玩家
func (g *Grid) GetPlayerIDs() []int {
	g.pIDLock.RLock()
	defer g.pIDLock.RUnlock()
	ids := make([]int, 0)
	for id, _ := range g.playerIDs {
		ids = append(ids, id)
	}
	return ids
}

// 打印信息方法
func (g *Grid) String() string {
	return fmt.Sprintf("Grid id: %d, minX:%d, maxX:%d, minY:%d, maxY:%d, playerIDs:%v",
		g.GID, g.MinX, g.MaxX, g.MinY, g.MaxY, g.playerIDs)
}
