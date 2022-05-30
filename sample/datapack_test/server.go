package main

import (
	"fmt"
	"io"
	"net"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xnet"
)

func main() {
	utils.InitConfig()

	listenner, err := net.Listen("tcp", "127.0.0.1:8777")
	if err != nil {
		fmt.Println("Listen fail err: ", err)
		return
	}

	fmt.Println("Start accepting connection. ..")
	for {
		conn, err := listenner.Accept()
		if err != nil {
			fmt.Println("server accept err: ", err)
			continue
		}

		go func(conn net.Conn) {
			dp := xnet.NewDataPack()

			for {
				headData := make([]byte, dp.GetHeadLen())

				_, err := io.ReadFull(conn, headData)
				if err != nil {
					fmt.Println("read head err: ", err)
					return
				}
				msgHead, err := dp.UnPack(headData)
				if err != nil {
					fmt.Println("unpack head fail: ", err)
					return
				}

				if msgHead.GetDataLen() > 0 {
					msg := msgHead.(*xnet.Message)
					msg.Data = make([]byte, msg.GetDataLen())

					// read bytes from io dps dataLen
					_, err := io.ReadFull(conn, msg.Data)
					if err != nil {
						fmt.Println("server read data fail: ", err)
						return
					}
					fmt.Println("===> Recv Msg ID: ", msg.Id, " len: ", msg.DataLen, ", data: ", string(msg.Data))
				}
			}
		}(conn)
	}
}
