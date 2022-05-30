package xnet

import (
	"fmt"

	"github.com/xeays/luffy/xiface"
)

type MsgHandler struct {
	Apis map[uint32]xiface.IRouter
}

func NewMsgHandler() *MsgHandler {
	return &MsgHandler{
		Apis: make(map[uint32]xiface.IRouter),
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
