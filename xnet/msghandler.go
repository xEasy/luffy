package xnet

import (
	"fmt"

	"github.com/xeays/luffy/utils"
	"github.com/xeays/luffy/xiface"
)

type MsgHandler struct {
	Apis           map[uint32]xiface.IRouter
	WorkerPoolSize uint32
	TaskQueue      []chan xiface.IRequest
}

func (mh *MsgHandler) StartWorkPool() {
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		// mallco space for current worker's queue
		mh.TaskQueue[i] = make(chan xiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)

		// start current worker, block until received new msg form taskqueue
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// SendMsgToTaskQueue send msg taks to taskqueue by request id using id/mod
func (mh *MsgHandler) SendMsgToTaskQueue(request xiface.IRequest) {
	// TODO worker pull task when free, all task in a stack queue
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID:", request.GetConnection().GetConnID(), ", requst msgID: ", request.GetMsgID(), " to workerID: ", workerID)

	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandler) StartOneWorker(workerID int, taskQueue chan xiface.IRequest) {
	fmt.Println("Worker ID: ", workerID, " is started.")
	for {
		select {
		case request := <-taskQueue:
			mh.DoMsghandler(request)
		}
	}
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
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		TaskQueue:      make([]chan xiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen),
	}
}
