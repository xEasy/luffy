// Package xnet provides tcp server
package xnet

import (
	"fmt"
	"net"

	"github.com/xeays/luffy/xiface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
}

func NewServer(name string) xiface.IServer {
	return &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
	}
}

func (s *Server) Start() {
	fmt.Printf("[Luffy] Server Listening at IP: %s, Port: %d Starging \n", s.IP, s.Port)
	go func() {
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("ResolveTCPAddr Fail with err: ", err)
			return
		}

		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("ListenTCP Fail with err: ", err)
		}

		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP err: ", err)
				continue
			}

			go func() {
				for {
					buf := make([]byte, 512)
					cnt, err := conn.Read(buf)
					if err != nil {
						fmt.Println("recv buf err ", err)
						continue
					}

					if _, err := conn.Write(buf[:cnt]); err != nil {
						fmt.Println("Write buf err: ", err)
						continue
					}
				}
			}()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("Server Stoping")
}

func (s *Server) Serve() {
	s.Start()
}
