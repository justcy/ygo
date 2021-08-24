package yiface

type IRouter interface {
	PreHandle(request IRequest) //在处理conn业务之前的hook方法
	Handle(request IRequest) //在处理conn业务之前的hook方法
	AfterHandle(request IRequest) //在处理conn业务之前的hook方法
}
