package znet

import (
	"errors"
	"fmt"
	"sync"
	"zinx/ziface"
)

type ConnManager struct {
	connMap  map[uint32]ziface.IConnection // 管理的连接信息
	connLock sync.RWMutex                  // 读写锁，保护连接map的并发
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connMap: make(map[uint32]ziface.IConnection),
	}
}

func (cm *ConnManager) AddConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	cm.connMap[conn.GetConnID()] = conn
	fmt.Printf("ConnID: %d added to ConnManager success, connNum: %d\n", conn.GetConnID(), cm.GetConnNum())
}

func (cm *ConnManager) DelConn(conn ziface.IConnection) {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()
	delete(cm.connMap, conn.GetConnID())
	fmt.Printf("ConnID: %d deleted from ConnManager success, connNum: %d\n", conn.GetConnID(), cm.GetConnNum())
}

func (cm *ConnManager) GetConnByID(connID uint32) (ziface.IConnection, error) {
	cm.connLock.RLock()
	defer cm.connLock.RUnlock()
	conn, ok := cm.connMap[connID]
	if !ok {
		return nil, errors.New(fmt.Sprintf("Connection with ID: %d not exists", connID))
	}
	return conn, nil
}

func (cm *ConnManager) GetConnNum() int {
	return len(cm.connMap)
}

func (cm *ConnManager) ClearConns() {
	cm.connLock.Lock()
	defer cm.connLock.Unlock()

	for connId, conn := range cm.connMap {
		// 停止
		conn.Stop()
		// 删除
		delete(cm.connMap, connId)
	}
	fmt.Printf("Clear all connections success, connNum: %v\n", len(cm.connMap))
}
