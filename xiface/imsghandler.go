package xiface

type IMsgHandler interface {
	DoMsghandler(request IRequest)
	AddRouter(msgID uint32, router IRouter)
	StartWorkPool()
	SendMsgToTaskQueue(request IRequest)
}
