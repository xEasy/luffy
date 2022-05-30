package xnet

type Message struct {
	Id      uint32
	Data    []byte
	DataLen uint32
}

func (m *Message) GetMsgId() uint32 {
	return m.Id
}

func (m *Message) GetData() []byte {
	return m.Data
}

func (m *Message) GetDataLen() uint32 {
	return m.DataLen
}

func (m *Message) SetData(data []byte) {
	m.Data = data
}

func (m *Message) SetMsgID(msgID uint32) {
	m.Id = msgID
}

func (m *Message) SetDataLen(dataLen uint32) {
	m.DataLen = dataLen
}
