package xiface

type IConManager interface {
	Add(conn IConnection)
	Remove(conn IConnection)
	Get(connID uint32) (IConnection, bool) // get connection by ID
	Len() int
	ClearConn()
}
