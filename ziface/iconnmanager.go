package ziface

type IConnManager interface {
	AddConn(conn IConnection)
	DelConn(conn IConnection) // TODO DelConnByID
	GetConnByID(connID uint32) (IConnection, error)
	GetConnNum() int
	ClearConns()
}
