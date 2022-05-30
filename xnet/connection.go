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
	HandleAPI    xiface.HandFunc
	ExitBuffChan chan bool
}

func NewConnection(conn *net.TCPConn, connID uint32, callback_api xiface.HandFunc) xiface.IConnection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		HandleAPI:    callback_api,
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
		cnt, err := conn.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err ", err)
			conn.ExitBuffChan <- true
			return
		}

		if err := conn.HandleAPI(conn.Conn, buf, cnt); err != nil {
			fmt.Println("connID ", conn.ConnID, " handle is error")
			conn.ExitBuffChan <- true
			return
		}
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
	return conn.RemoteAddr()
}
