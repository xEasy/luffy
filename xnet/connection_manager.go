package xnet

import (
	"fmt"
	"sync"

	"github.com/xeays/luffy/xiface"
)

type ConnManager struct {
	connections map[uint32]xiface.IConnection
	connLock    sync.RWMutex
}

func (c *ConnManager) Add(conn xiface.IConnection) {
	c.connLock.Lock()
	c.connections[conn.GetConnID()] = conn
	c.connLock.Unlock()

	fmt.Println("Connection ID: ", conn.GetConnID(), " add to ConnManager successfully, conn Len: ", c.Len())
}

func (c *ConnManager) Remove(conn xiface.IConnection) {
	c.connLock.Lock()
	delete(c.connections, conn.GetConnID())
	c.connLock.Unlock()

	fmt.Println("Connection ID: ", conn.GetConnID(), " remove from ConnManager successfully, conn Len: ", c.Len())
}

func (c *ConnManager) Get(connID uint32) (xiface.IConnection, bool) {
	c.connLock.RLock()
	conn, ok := c.connections[connID]
	c.connLock.RUnlock()

	return conn, ok
}

func (c *ConnManager) Len() int {
	return len(c.connections)
}

func (c *ConnManager) ClearConn() {
	c.connLock.Lock()
	for connID, conn := range c.connections {
		conn.Stop()
		delete(c.connections, connID)
	}
	c.connLock.Unlock()
	fmt.Println("Clear All Connections succ! conn len: ", c.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]xiface.IConnection),
	}
}
