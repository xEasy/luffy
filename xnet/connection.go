package xnet

import (
	"fmt"
	"net"

	"github.com/xeays/luffy/xiface"
)

type Connection struct {
	ConnID       uint32
	Conn         *net.TCPConn
	isClosed     bool
	Router       xiface.IRouter
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, router xiface.IRouter) xiface.IConnection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader is running id: ", conn.ConnID)
	defer fmt.Println(conn.RemoteAddr().String(), " conn reader exit!")
	defer conn.Stop()

	for {
		buf := make([]byte, 512)
		_, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			conn.ExitBuffChan <- true
			return
		}

		req := &Request{
			data: buf,
			conn: conn,
		}

		go func(request xiface.IRequest) {
			conn.Router.Handle(request)
		}(req)
	}
}

// Start connection to work
func (conn *Connection) Start() {

	// create goroutinue to handle connection
	go conn.StartReader()

	for {
		select {
		case <-conn.ExitBuffChan:
			// recv exit signal and stop block
			return
		}
	}
}

func (conn *Connection) Stop() {
	if conn.isClosed == true {
		return
	}

	conn.isClosed = true

	// TODO run regieste stop callback func

	// close tcp socket
	conn.Conn.Close()

	// notify exit message subscriber
	conn.ExitBuffChan <- true

	// close exit channel
	close(conn.ExitBuffChan)
}

func (conn *Connection) GetTCPConnection() *net.TCPConn {
	return conn.Conn
}

func (conn *Connection) GetConnID() uint32 {
	return conn.ConnID
}

func (conn *Connection) RemoteAddr() net.Addr {
	return conn.Conn.RemoteAddr()
}
