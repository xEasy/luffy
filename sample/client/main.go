package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/xeays/luffy/xnet"
)

func main() {
	fmt.Println("Client starting")
	conn, err := net.Dial("tcp4", "0.0.0.0:8777")
	if err != nil {
		fmt.Println("client dial fail err: ", err)
		return
	}

	for {
		dp := xnet.NewDataPack()

		msg, _ := dp.Pack(xnet.NewMsgPackage(1, []byte("hello luffy server")))

		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("writ err : ", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head fail")
			break
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err: ", err)
		}

		if msgHead.GetDataLen() > 0 {
			// read msg data
			msg := msgHead.(*xnet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// read bytes from io by dataLen
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data fail: ", err)
				return
			}

			fmt.Println("====> Recv Msg ID: ", msg.Id, ", Len: ", msg.DataLen, ", data: ", string(msg.Data))
		}

		time.Sleep(time.Second)
	}
}
