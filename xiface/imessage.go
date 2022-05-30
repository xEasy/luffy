package xiface

type IMessage interface {
	GetDataLen() uint32
	GetMsgId() uint32
	GetData() []byte

	SetMsgID(uint32)
	SetData([]byte)
	SetDataLen(uint32)
}
