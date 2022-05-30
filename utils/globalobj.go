package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/xeays/luffy/xiface"
)

type GlobalObj struct {
	TcpServer xiface.IServer // current global server obj
	Host      string         // the host that server hosting
	TcpPort   int            // the port that server lintening
	Name      string         // current server name
	Version   string         // current luffy version

	MaxPacketSize uint32 // max data package
	MaxConn       int    // max connections on current server

	PoolSize         uint32 // msg pool size
	WorkerPoolSize   uint32 // each msgPool's max worker pool size
	MaxWorkerTaskLen uint32 // task worker queue's max len

	WriterPoolSize   uint32
	MaxWriterTaskLen uint32
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/luffy.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func InitConfig() {
	GlobalObject = &GlobalObj{
		Host:    "0.0.0.0",
		TcpPort: 8999,
		Name:    "OPice",
		Version: "v0.1",

		MaxPacketSize: 4096,
		MaxConn:       12000,

		PoolSize:         10,
		WorkerPoolSize:   100,
		MaxWorkerTaskLen: 1024,

		WriterPoolSize:   100,
		MaxWriterTaskLen: 1024,
	}

	GlobalObject.Reload()
}
