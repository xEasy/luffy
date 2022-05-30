package xiface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error

	// set propertity value
	SetProperty(key string, value any)
	// get propertity value
	GetProperty(key string) (any, bool)
	// remove connection's property
	RemoveProperty(key string)
}
