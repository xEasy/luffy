## Luffy

Luffy is an light weight TCP server, rewrite inspire by [aceld/znix](https://github.com/aceld/zinx)

### Feature

- mutil router
- TLV message model
- separete read & write
- message router & handling worker pool 

### NewFeature

- flexible workerpool
- writer pool
- consistent hash

### Usage

#### server

```golang
package main

import (
	"fmt"

	"github.com/xeays/luffy/xiface"
	"github.com/xeays/luffy/xnet"
)

// can be any struct that realize xiface.IRouter
type PingRouter struct {
	xnet.BaseRouter
}

func main() {
    // make a new TCP server
	s := xnet.NewServer("Luffy 0.1")

    // start serving with on block
	s.Serve()

    // config message router
	s.AddRouter(0, &PingRouter{})

	select {}
}

func (r *PingRouter) Handle(request xiface.IRequest) {
	fmt.Println("PingRouter is called")

    // send message to client
	err := request.GetConnection().SendMsg(0, []byte("ping.. ping.. pong.. pong.."))
	if err != nil {
		fmt.Println("call back PingRouter err", err)
	}
}
```


#### client

```golang
package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/xeays/luffy/xnet"
)

func main() {
	fmt.Println("Client starting")
	conn, err := net.Dial("tcp4", "0.0.0.0:8777")
	if err != nil {
		fmt.Println("client dial fail err: ", err)
		return
	}

	for {
		dp := xnet.NewDataPack()

		msg, _ := dp.Pack(xnet.NewMsgPackage(1, []byte("hello luffy server")))

		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("writ err : ", err)
			return
		}

		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData)
		if err != nil {
			fmt.Println("read head fail")
			break
		}

		msgHead, err := dp.UnPack(headData)
		if err != nil {
			fmt.Println("server unpack err: ", err)
		}

		if msgHead.GetDataLen() > 0 {
			// read msg data
			msg := msgHead.(*xnet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			// read bytes from io by dataLen
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("read msg data fail: ", err)
				return
			}

			fmt.Println("====> Recv Msg ID: ", msg.Id, ", Len: ", msg.DataLen, ", data: ", string(msg.Data))
		}

		time.Sleep(time.Microsecond * 100)
	}
}
```

#### config

- config file path: `conf/luffy.json`
- example:

```json
{
  "Name": "Demo server",
  "Host": "0.0.0.0",
  "TcpPort": 8777,
  "MaxConn": 30
	"MaxPacketSize": 1024, // max data package
	MaxCon: 65555,    // max connections on current server

	PoolSize: 10, // msgPools size
	WorkerPoolSize: 100 // each msgPool's max worker pool size
	MaxWorkerTaskLen: 1024 // each msgPool's task queue buf size

	WriterPoolSize: 10, // writer pool's worker size
	MaxWriterTaskLen: 1024, // writer pool's task queue buf size
}
```

#### workerpool

```golang
    // make a new workerpool
	pool := workerpool.NewWorkPool("testBatch", 5, 10)

    // start dispatching job with no block
	pool.Start()

    // let pool to init complete
	time.Sleep(time.Microsecond * 100) 

    // an example job:
	// inc counter with step 2
    counter := 0
	for i := 0; i < 100; i++ {
		job := func(args ...any) {
            // run any func you want
			prime := args[0].(int32)
			atomic.CompareAndSwapInt32(&counter, counter, counter+prime)
		}

		pool.Enqueue(job, int32(2))
	}

    // call release if shutting down
	// release will block until all jobs done
	pool.Release()

    // no block release
    go pool.Release()
    select {
        case: <- pool.Done()
            // this code will run until pool release done
            .....
    }
```

#### consistent hash

```golang
    // make a consistent hash obj
	cs := NewConsistent()

    // add nodes with  virtual nodes size
	cs.Add("one", 150)
	cs.Add("two", 100)
	err := cs.Add("three", 100)
	if err == nil {
		//...
	}

    // get node by <key>
	node := cs.GetNode("hello")
    if node == "" {
        // not found node
    }

    // remove node
	err = cs.Remove("one")
    if err != nil {
        // Oooops...
    }
```
