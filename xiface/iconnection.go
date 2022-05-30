package xiface

import "net"

type IConnection interface {
	Start()
	Stop()
	GetTCPConnection() *net.TCPConn
	GetConnID() uint32
	RemoteAddr() net.Addr
	SendMsg(msgID uint32, data []byte) error
}

type HandFunc func(*net.TCPConn, []byte, int) error
