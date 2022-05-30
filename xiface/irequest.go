package xiface

type IRequest interface {
	GetData() []byte
	GetTCPConnection() IConnection
	GetMsgID() uint32
}
