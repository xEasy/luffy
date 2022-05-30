package main

import (
	"fmt"
	"net"

	"github.com/xeays/luffy/xnet"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:8777")
	if err != nil {
		fmt.Println("client dial fail: ", err)
		return
	}

	dp := xnet.NewDataPack()

	msg1 := &xnet.Message{
		Id:      0,
		DataLen: 5,
		Data:    []byte("hello"),
	}

	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("packa msg1 err: ", err)
		return
	}

	msg2 := &xnet.Message{
		Id:      0,
		DataLen: 7,
		Data:    []byte("world!!"),
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("pack msg2 fail: ", err)
		return
	}

	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)

	select {}
}
