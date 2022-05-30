package xnet

import "github.com/xeays/luffy/xiface"

type Request struct {
	id   string
	conn xiface.IConnection
	msg  xiface.IMessage
}

func (r *Request) GetID() string {
	return r.id
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}

func (r *Request) GetConnection() xiface.IConnection {
	return r.conn
}
