package xnet

import (
	"fmt"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/workerpool"
	"github.com/xeays/luffy/xiface"
)

type MsgHandler struct {
	Apis           map[uint32]xiface.IRouter
	PoolSize       uint32
	WorkerPoolSize uint32
	MsgPools       []*workerpool.Pool
}

func (mh *MsgHandler) StartWorkPool() {
	if mh.MsgPools == nil {
		mh.MsgPools = make([]*workerpool.Pool, mh.WorkerPoolSize)
	}

	for i := 0; i < int(mh.PoolSize); i++ {
		poolName := fmt.Sprintf("LuffyMsgPool:%d", i)
		pool := workerpool.NewWorkPool(poolName, mh.WorkerPoolSize, utils.GlobalObject.MaxWorkerTaskLen)
		mh.MsgPools[i] = pool
		pool.Start()
	}
}

// SendMsgToTaskQueue send msg taks to taskqueue by request id using id/mod
func (mh *MsgHandler) SendMsgToTaskQueue(request xiface.IRequest) {
	// each pool's worker pull task when free, all task in job queue
	poolID := request.GetConnection().GetConnID() % mh.PoolSize

	fmt.Println("Add ConnID:", request.GetConnection().GetConnID(), ", requst msgID: ", request.GetMsgID(), " to workerID: ", poolID)

	pool := mh.MsgPools[poolID]
	pool.Enqueue(func(args ...any) {
		mh := args[0].(*MsgHandler)
		request := args[1].(xiface.IRequest)
		mh.DoMsghandler(request)
	}, mh, request)
}

func (mh *MsgHandler) DoMsghandler(request xiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgID()]
	if !ok {
		fmt.Println("api msgID: ", request.GetMsgID(), " is not FOUND")
		return
	}

	handler.Handle(request)
}

func (mh *MsgHandler) AddRouter(msgID uint32, router xiface.IRouter) {
	if _, ok := mh.Apis[msgID]; ok {
		panic(fmt.Sprintf("repeated api, msgID: %d", msgID))
	}

	mh.Apis[msgID] = router
	fmt.Println("Add api msgID : ", msgID)
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis:           make(map[uint32]xiface.IRouter),
		PoolSize:       utils.GlobalObject.PoolSize,
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		MsgPools:       nil,
	}
}
