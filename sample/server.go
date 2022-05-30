package main

import (
	"fmt"

	"github.com/xeays/luffy/xiface"
	"github.com/xeays/luffy/xnet"
)

type PingRouter struct {
	xnet.BaseRouter
}

type BoomRouter struct {
	xnet.BaseRouter
}

func main() {
	s := xnet.NewServer("Luffy 0.1")
	s.Serve()
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &BoomRouter{})

	select {}
}

func (r *PingRouter) Handle(request xiface.IRequest) {
	fmt.Println("PingRouter is called")

	err := request.GetConnection().SendMsg(0, []byte("ping.. ping.. pong.. pong.."))
	if err != nil {
		fmt.Println("call back PingRouter err", err)
	}
}

func (r *BoomRouter) Handle(req xiface.IRequest) {
	fmt.Println("PingRouter is called")

	err := req.GetConnection().SendMsg(1, []byte("boom boom boom boom.."))
	if err != nil {
		fmt.Println("call back PongRouter err", err)
	}
}
