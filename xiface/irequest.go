package xiface

type IRequest interface {
	GetData() []byte
	GetTCPConnection() IConnection
}
