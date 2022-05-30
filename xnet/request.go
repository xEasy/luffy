package xnet

import "github.com/xeays/luffy/xiface"

type Request struct {
	conn xiface.IConnection
	msg  xiface.IMessage
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetTCPConnection() xiface.IConnection {
	return r.conn
}
