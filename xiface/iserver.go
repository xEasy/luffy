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
}
