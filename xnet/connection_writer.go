package xnet

import (
	"fmt"
	"sync"

	"github.com/xeays/luffy/workerpool"
	"github.com/xeays/luffy/xiface"
)

type ConnectionWriter struct {
	mu            sync.RWMutex
	writerSize    uint32
	writerBufSize uint32
	stoped        bool
	pool          *workerpool.Pool
}

func NewConnectionWriter(wSize uint32, bufSize uint32) xiface.IConnWriter {
	return &ConnectionWriter{
		writerSize:    wSize,
		writerBufSize: bufSize,
		stoped:        false,
		pool:          nil,
	}
}

func (cw *ConnectionWriter) Start() {
	if cw.pool != nil {
		return
	}
	pool := workerpool.NewWorkPool("LuffyWriter", cw.writerSize, cw.writerBufSize)
	cw.pool = pool
	pool.Start()
}

func (cw *ConnectionWriter) Stop() {
	if cw.pool == nil {
		return
	}
	cw.mu.Lock()
	if cw.stoped {
		cw.mu.Unlock()
		return
	}
	cw.stoped = true
	cw.mu.Unlock()

	cw.pool.Release()
}

func (cw *ConnectionWriter) Write(c xiface.IConnection, data []byte) {
	cw.mu.RLock()

	if cw.stoped {
		cw.mu.RUnlock()
		return
	}

	cw.pool.Enqueue(func(args ...any) {
		c := args[0].(xiface.IConnection)
		data := args[1].([]byte)
		if _, err := c.GetTCPConnection().Write(data); err != nil {
			fmt.Println("Send data error, err: ", err, " RemoteAddr: ", c.GetTCPConnection().RemoteAddr())
		}
	}, c, data)

	cw.mu.RUnlock()
}
