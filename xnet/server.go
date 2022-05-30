// Package xnet provides tcp server
package xnet

import (
	"fmt"
	"net"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xiface"
)

type Server struct {
	Name        string
	IPVersion   string
	IP          string
	Port        int
	MsgHandler  xiface.IMsgHandler
	connManager xiface.IConManager
}

func NewServer(name string) xiface.IServer {

	utils.InitConfig()
	utils.GlobalObject.Reload()

	return &Server{
		Name:        utils.GlobalObject.Name,
		IPVersion:   "tcp4",
		IP:          utils.GlobalObject.Host,
		Port:        utils.GlobalObject.TcpPort,
		MsgHandler:  NewMsgHandler(),
		connManager: NewConnManager(),
	}
}

func (s *Server) Start() {
	fmt.Printf("[Luffy] Server Listening at IP: %s, Port: %d Starging \n", s.IP, s.Port)
	fmt.Printf("[Luffy] Version: %s, MaxConn: %d, MaxPacketSize: %d \n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize,
	)
	go func() {
		// 1 start msgHandler workPool
		s.MsgHandler.StartWorkPool()

		// 2 get a tcp addr
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

			// close conn if cid greater than maxID
			if s.connManager.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			dealConn := NewConnection(s, conn, cid, s.MsgHandler)

			cid++

			go dealConn.Start()
		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("Server Stoping ...")
	s.connManager.ClearConn()
	fmt.Println("Server Stoped")
}

func (s *Server) Serve() {
	s.Start()
}

func (s *Server) AddRouter(msgID uint32, router xiface.IRouter) {
	s.MsgHandler.AddRouter(msgID, router)
	fmt.Println("AddRouter succ! ")
}

// Get connection manager
func (s *Server) GetConnMgr() xiface.IConManager {
	return s.connManager
}
