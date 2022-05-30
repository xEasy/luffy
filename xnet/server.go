// Package xnet provides tcp server
package xnet

import (
	"errors"
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

func callbackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	fmt.Println("[conn handler] callbackToClient ...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf error: ", err)
		return errors.New("callbackToClient error")
	}
	return nil
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

		// TODO id generate with func
		var cid uint32
		cid = 0

		// start server listen
		for {
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("AcceptTCP err: ", err)
				continue
			}

			// TODO close conn if cid greater than maxID

			dealConn := NewConnection(conn, cid, callbackToClient)

			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("Server Stoping")
}

func (s *Server) Serve() {
	s.Start()
}
