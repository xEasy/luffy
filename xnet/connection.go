package xnet

import (
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xiface"
)

type Connection struct {
	Server       xiface.IServer
	ConnID       uint32
	Conn         *net.TCPConn
	isClosed     bool
	ExitBuffChan chan bool
	MsgHandler   xiface.IMsgHandler
	msgChan      chan []byte

	// connection properties
	properties map[string]any
	ppLock     sync.RWMutex
}

// set propertity value
func (c *Connection) SetProperty(key string, value any) {
	c.ppLock.Lock()
	c.properties[key] = value
	c.ppLock.Unlock()
}

// get propertity value
func (c *Connection) GetProperty(key string) (any, bool) {
	c.ppLock.RLock()
	value, ok := c.properties[key]
	c.ppLock.RUnlock()
	return value, ok
}

// remove connection's property
func (c *Connection) RemoveProperty(key string) {
	c.ppLock.Lock()
	delete(c.properties, key)
	c.ppLock.Unlock()
}

func NewConnection(server xiface.IServer, conn *net.TCPConn, connID uint32, msgHandler xiface.IMsgHandler) xiface.IConnection {
	c := &Connection{
		Server:       server,
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		MsgHandler:   msgHandler,
		ExitBuffChan: make(chan bool, 1),
		msgChan:      make(chan []byte),
	}

	return c
}

func (c *Connection) StartWriter() {
	fmt.Println("[Writer Goroutine is running...]")
	defer fmt.Println(c.RemoteAddr().String(), " [conn Writer exit!]")

	for {
		select {
		case data := <-c.msgChan:
			if _, err := c.Conn.Write(data); err != nil {
				fmt.Println("Send data err: ", err, " Conn Writer exit!")
				return
			}
		case <-c.ExitBuffChan:
			return
		}
	}
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader is running id: ", conn.ConnID)
	defer fmt.Println(conn.RemoteAddr().String(), " conn reader exit!")
	defer conn.Stop()

	for {
		dp := NewDataPackWithMaxSize(utils.GlobalObject.MaxPacketSize)

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

		if utils.GlobalObject.WorkerPoolSize > 0 {
			conn.MsgHandler.SendMsgToTaskQueue(&req)
		} else {
			go conn.MsgHandler.DoMsghandler(&req)
		}
	}
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	dp := NewDataPackWithMaxSize(utils.GlobalObject.MaxPacketSize)
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack message fail: ", err)
		return errors.New("pack error msg")
	}

	// write to client throught msgChan
	c.msgChan <- msg

	return nil
}

// Start connection to work
func (conn *Connection) Start() {

	// create goroutinue to handle connection
	go conn.StartReader()
	// create goroutinue to write
	go conn.StartWriter()

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

	// remove conn from connManger
	conn.Server.GetConnMgr().Remove(conn)

	// close conn's channel
	close(conn.ExitBuffChan)
	close(conn.msgChan)
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
