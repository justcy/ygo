package main

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ynet"
)

//ping test 自定义路由
type PingRouter struct {
	ynet.BaseRouter//一定要先基础BaseRouter
}

//Test PreHandle
func (this *PingRouter) PreHandle(request yiface.IRequest) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping ....\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}
//Test Handle
func (this *PingRouter) Handle(request yiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}

//Test PostHandle
func (this *PingRouter) AfterHandle(request yiface.IRequest) {
	fmt.Println("Call Router PostHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping .....\n"))
	if err !=nil {
		fmt.Println("call back ping ping ping error")
	}
}

func main() {

	//1 创建一个server 句柄 s
	s := ynet.NewServer()

	s.AddRouter(&PingRouter{})
	//2 开启服务
	s.Server()
}
