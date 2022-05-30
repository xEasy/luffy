package xiface

type IRequest interface {
	GetData() []byte
	GetConnection() IConnection
	GetMsgID() uint32
	GetID() string
}
