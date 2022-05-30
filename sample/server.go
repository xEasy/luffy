package main

import "github.com/xeays/luffy/xnet"

func main() {
	s := xnet.NewServer("Luffy 0.1")
	s.Serve()

	select {}
}
