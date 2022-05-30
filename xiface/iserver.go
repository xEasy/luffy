// Package xiface provides LuffyTcpServer
package xiface

type IServer interface {
	Start()
	Stop()
	Serve()
}
