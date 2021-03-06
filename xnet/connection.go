package xnet

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"

	"github.com/google/uuid"
	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xiface"
)

type Connection struct {
	mu         sync.Mutex
	Server     xiface.IServer
	ConnID     uint32
	Conn       *net.TCPConn
	isClosed   bool
	MsgHandler xiface.IMsgHandler

	ctx    context.Context
	cancel context.CancelFunc

	// connection properties
	properties map[string]any
	ppLock     sync.RWMutex
}

func NewConnection(server xiface.IServer, conn *net.TCPConn, connID uint32, msgHandler xiface.IMsgHandler) xiface.IConnection {
	c := &Connection{
		Server:     server,
		Conn:       conn,
		ConnID:     connID,
		isClosed:   false,
		MsgHandler: msgHandler,
		properties: nil,
	}

	// add client connection to connManager
	c.Server.GetConnMgr().Add(c)

	return c
}

func (conn *Connection) StartReader() {
	fmt.Println("Reader is running id: ", conn.ConnID)
	defer fmt.Println(conn.RemoteAddr().String(), " conn reader exit!")
	defer conn.Stop()

	for {
		select {
		case <-conn.ctx.Done():
			return
		default:
			if err := conn.read(); err != nil {
				// quit when read err
				return
			}
		}
	}
}

func (conn *Connection) read() error {
	dp := NewDataPackWithMaxSize(utils.GlobalObject.MaxPacketSize)

	// read client msg head
	headData := make([]byte, dp.GetHeadLen())
	if _, err := io.ReadFull(conn.GetTCPConnection(), headData); err != nil {
		fmt.Println("read msg head err: ", err)
		return err
	}

	msg, err := dp.UnPack(headData)
	if err != nil {
		fmt.Println("unpack msg head err: ", err)
		return err
	}

	var msgData []byte
	if msg.GetDataLen() > 0 {
		msgData = make([]byte, msg.GetDataLen())
		if _, err := io.ReadFull(conn.GetTCPConnection(), msgData); err != nil {
			fmt.Println("read msg data fail: ", err)
			return err
		}
	}

	msg.SetData(msgData)

	var reqID string
	reqUUID, err := uuid.NewRandom()
	if err != nil {
		reqID = "fakeID"
	} else {
		reqID = reqUUID.String()
	}
	req := Request{
		id:   reqID,
		msg:  msg,
		conn: conn,
	}

	if utils.GlobalObject.WorkerPoolSize > 0 {
		conn.MsgHandler.SendMsgToTaskQueue(&req)
	} else {
		go conn.MsgHandler.DoMsghandler(&req)
	}
	return nil
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	c.mu.Lock()
	if c.isClosed == true {
		c.mu.Unlock()
		return errors.New("Connection closed when send msg")
	}
	c.mu.Unlock()

	dp := NewDataPackWithMaxSize(utils.GlobalObject.MaxPacketSize)
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("pack message fail: ", err)
		return errors.New("pack error msg")
	}

	// write to client throught writer pool
	c.Server.GetConnWriter().Write(c, msg)

	return nil
}

// Start connection to work
func (conn *Connection) Start() {
	conn.ctx, conn.cancel = context.WithCancel(context.Background())
	// create goroutinue to handle connection
	go conn.StartReader()

	// client connection start callback
	conn.Server.CallOnConnStart(conn)

	select {
	case <-conn.ctx.Done():
		// remove conn from connManager
		conn.Server.GetConnMgr().Remove(conn)
		return
	}
}

func (conn *Connection) Stop() {
	conn.mu.Lock()

	if conn.isClosed == true {
		conn.mu.Unlock()
		return
	}

	conn.isClosed = true

	// close tcp socket
	conn.Conn.Close()

	// notify exit message subscriber
	conn.cancel()

	// unlock until close action done
	conn.mu.Unlock()

	// run regieste stop callback func
	conn.Server.CallOnConnStop(conn)
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

// set propertity value
func (c *Connection) SetProperty(key string, value any) {
	c.ppLock.Lock()
	if c.properties == nil {
		c.properties = make(map[string]any)
	}
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
