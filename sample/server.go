package main

import (
	"fmt"

	"github.com/xeays/luffy/xiface"
	"github.com/xeays/luffy/xnet"
)

type PingRouter struct {
	xnet.BaseRouter
}

func main() {
	s := xnet.NewServer("Luffy 0.1")
	s.Serve()
	s.AddRouter(&PingRouter{})

	select {}
}

func (r *PingRouter) Handle(request xiface.IRequest) {
	fmt.Println("PingRouter is called")

	err := request.GetTCPConnection().SendMsg(1, []byte("ping.. ping.. pong.. pong.."))
	if err != nil {
		fmt.Println("call back PingRouter err", err)
	}
}
