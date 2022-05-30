package xnet

import "github.com/xeays/luffy/xiface"

type Request struct {
	conn xiface.IConnection
	data []byte
}

func (r *Request) GetData() []byte {
	return r.data
}

func (r *Request) GetTCPConnection() xiface.IConnection {
	return r.conn
}
