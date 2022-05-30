package xiface

type IRouter interface {
	Handle(req IRequest)
}
