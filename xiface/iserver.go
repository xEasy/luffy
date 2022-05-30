// Package xiface provides LuffyTcpServer
package xiface

type IServer interface {
	// bootup the server
	Start()

	// stop the server
	Stop()

	// let server start serving
	Serve()

	// router func
	AddRouter(msgID uint32, router IRouter)

	// Get connection manager
	GetConnMgr() IConManager

	// Get connection Writer
	GetConnWriter() IConnWriter

	// SetOnConnStart set connection start callback
	SetOnConnStart(func(IConnection))

	// SetOnConnStop set connection stop callback
	SetOnConnStop(func(IConnection))

	//CallOnConnStart run conn start callback
	CallOnConnStart(conn IConnection)

	//CallOnConnStop run conn stop callback
	CallOnConnStop(conn IConnection)
}
