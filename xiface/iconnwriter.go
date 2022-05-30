package xiface

type IConnWriter interface {
	Start()
	Stop()
	Write(c IConnection, data []byte)
}
