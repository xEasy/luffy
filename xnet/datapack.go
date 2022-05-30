package xnet

import (
	"bytes"
	"encoding/binary"
	"errors"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xiface"
)

type DataPack struct{}

func (dp *DataPack) GetHeadLen() uint32 {
	// Id uint32(4 bytes) + DataLen uint32(4 bytes)
	return 8
}

func (dp *DataPack) Pack(msg xiface.IMessage) ([]byte, error) {
	dataBuff := bytes.NewBuffer([]byte{})

	// write dataLen using LittleEndian mode
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}

	// write Id
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}

	// write data
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}

	return dataBuff.Bytes(), nil
}

func (dp *DataPack) UnPack(binaryData []byte) (xiface.IMessage, error) {
	dataBuff := bytes.NewReader(binaryData)
	msg := &Message{}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}

	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}

	if utils.GlobalObject.MaxPacketSize > 0 && utils.GlobalObject.MaxPacketSize < msg.DataLen {
		return nil, errors.New("Too large msg data recv!")
	}

	// return msg head, read conn's data using dataLen storing in msg.head
	return msg, nil
}
