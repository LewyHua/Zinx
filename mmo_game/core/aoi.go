package core

import "fmt"

const (
	AOI_MIN_X  int = 85
	AOI_MAX_X  int = 410
	AOI_CNTS_X int = 10
	AOI_MIN_Y  int = 75
	AOI_MAX_Y  int = 400
	AOI_CNTS_Y int = 20
)

// AOIManager AOI 区域管理模块
type AOIManager struct {
	MinX      int           //区域左边界坐标
	MaxX      int           //区域右边界坐标
	XGridNums int           //x方向格子的数量
	MinY      int           //区域上边界坐标
	MaxY      int           //区域下边界坐标
	YGridNums int           //y方向的格子数量
	grids     map[int]*Grid //当前区域中都有哪些格子，key=格子ID， value=格子对象
}

func NewAOIManager(minX, maxX, xGridNums, minY, maxY, yGridNums int) *AOIManager {
	aoiMgr := &AOIManager{
		MinX:      minX,
		MaxX:      maxX,
		XGridNums: xGridNums,
		MinY:      minY,
		MaxY:      maxY,
		YGridNums: yGridNums,
		grids:     make(map[int]*Grid),
	}

	// 初始化区域内所有grid
	for y := 0; y < yGridNums; y++ {
		for x := 0; x < xGridNums; x++ {
			// grid id
			gid := y + x*yGridNums
			aoiMgr.grids[gid] = NewGrid(gid,
				aoiMgr.MinX+x*aoiMgr.gridWidth(),
				aoiMgr.MinX+(x+1)*aoiMgr.gridWidth(),
				aoiMgr.MinY+y*aoiMgr.gridLength(),
				aoiMgr.MinY+(y+1)*aoiMgr.gridLength())
		}
	}
	return aoiMgr
}

// 得到每个格子在x轴方向的宽度
func (m *AOIManager) gridWidth() int {
	return (m.MaxX - m.MinX) / m.XGridNums
}

// 得到每个格子在x轴方向的长度
func (m *AOIManager) gridLength() int {
	return (m.MaxY - m.MinY) / m.YGridNums
}

// 打印信息方法
func (m *AOIManager) String() string {
	s := fmt.Sprintf("AOIManagr:\nminX:%d, maxX:%d, xGridNums:%d, minY:%d, maxY:%d, yGridNums:%d\n Grids in AOI Manager:\n",
		m.MinX, m.MaxX, m.XGridNums, m.MinY, m.MaxY, m.YGridNums)
	for _, grid := range m.grids {
		s += fmt.Sprintln(grid)
	}

	return s
}

// GetSurroundingGridsByGID 根据格子的gID得到当前周边的九宫格信息
func (m *AOIManager) GetSurroundingGridsByGID(gID int) (grids []*Grid) {
	//判断gID是否存在
	if _, ok := m.grids[gID]; !ok {
		return
	}

	//将当前gid添加到九宫格中
	grids = append(grids, m.grids[gID])

	//根据gid得到当前格子所在的X轴编号
	idx := gID % m.XGridNums

	//判断当前idx左边是否还有格子
	if idx > 0 {
		grids = append(grids, m.grids[gID-1])
	}
	//判断当前的idx右边是否还有格子
	if idx < m.XGridNums-1 {
		grids = append(grids, m.grids[gID+1])
	}

	//将x轴当前的格子都取出，进行遍历，再分别得到每个格子的上下是否有格子

	//得到当前x轴的格子id集合
	gidsX := make([]int, 0, len(grids))
	for _, v := range grids {
		gidsX = append(gidsX, v.GID)
	}

	//遍历x轴格子
	for _, v := range gidsX {
		//计算该格子处于第几列
		idy := v / m.XGridNums

		//判断当前的idy上边是否还有格子
		if idy > 0 {
			grids = append(grids, m.grids[v-m.XGridNums])
		}
		//判断当前的idy下边是否还有格子
		if idy < m.YGridNums-1 {
			grids = append(grids, m.grids[v+m.XGridNums])
		}
	}
	return
}

// GetGIDByPos 通过坐标获取对应的格子ID
func (m *AOIManager) GetGIDByPos(x, y float32) int {
	gx := (int(x) - m.MinX) / m.gridWidth()
	gy := (int(y) - m.MinY) / m.gridLength()
	return gy*m.XGridNums + gx
}

// GetPIDsByPos 通过坐标得到周边九宫格内的全部PlayerIDs
func (m *AOIManager) GetPIDsByPos(x, y float32) (playerIDs []int) {
	//根据横纵坐标得到当前坐标属于哪个格子ID
	gID := m.GetGIDByPos(x, y)

	//根据格子ID得到周边九宫格的信息
	grids := m.GetSurroundingGridsByGID(gID)
	for _, grid := range grids {
		playerIDs = append(playerIDs, grid.GetPlayerIDs()...)
		//fmt.Printf("===> grid ID : %d, pids : %v  ====", grid.GID, grid.GetPlayerIDs())
	}

	return
}

// GetPlayerIDsByGID 通过GID获取当前格子的全部playerID
func (m *AOIManager) GetPlayerIDsByGID(gID int) (playerIDs []int) {
	playerIDs = m.grids[gID].GetPlayerIDs()
	return
}

// DelPIDFromGrid 移除一个格子中的PlayerID
func (m *AOIManager) DelPIDFromGrid(pID, gID int) {
	m.grids[gID].Del(pID)
}

// AddPIDToGrid 添加一个PlayerID到一个格子中
func (m *AOIManager) AddPIDToGrid(pID, gID int) {
	m.grids[gID].Add(pID)
}

// AddToGridByPos 通过横纵坐标添加一个Player到一个格子中
func (m *AOIManager) AddToGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Add(pID)
}

// DelFromGridByPos 通过横纵坐标把一个Player从对应的格子中删除
func (m *AOIManager) DelFromGridByPos(pID int, x, y float32) {
	gID := m.GetGIDByPos(x, y)
	grid := m.grids[gID]
	grid.Del(pID)
}
