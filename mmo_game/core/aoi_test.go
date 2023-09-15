package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	fmt.Println(aoiMgr)
}

// 0 0 0 0 0 0 0 1
// 0 0 0 0 0 0 1 0

func TestAOIManagerSurroundGridsByGID(t *testing.T) {
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)
	grids := aoiMgr.GetSurroundingGridsByGID(0)
	for _, grid := range grids {
		fmt.Println(grid)
	}
}
