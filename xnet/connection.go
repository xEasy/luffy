package xnet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/xeays/luffy/xiface"
)

type Connection struct {
	ConnID       uint32
	Conn         *net.TCPConn
	isClosed     bool
	ExitBuffChan chan bool
	MsgHandler   xiface.IMsgHandler
}

func NewConnection(conn *net.TCPConn, connID uint32, msgHandler xiface.IMsgHandler) xiface.IConnection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader is running id: ", conn.ConnID)
	defer fmt.Println(conn.RemoteAddr().String(), " conn reader exit!")
	defer conn.Stop()

	for {
		dp := NewDataPack()

		// read client msg head
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
			fmt.Println("read msg head err: ", err)
			conn.ExitBuffChan <- true
			continue
		}

		msg, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("unpack msg head err: ", err)
		}

		var msgData []byte
		if msg.GetDataLen() > 0 {
			msgData = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(conn.GetTCPConnection(), msgData); err != nil {
				fmt.Println("read msg data fail: ", err)
				conn.ExitBuffChan <- true
				continue
			}
		}

		msg.SetData(msgData)

		req := Request{
			msg:  msg,
			conn: conn,
		}

		go conn.MsgHandler.DoMsghandler(&req)
	}
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack message fail: ", err)
		return errors.New("pack error msg")
	}

	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("write msg id ", msgId, " error ")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}
	return nil
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
