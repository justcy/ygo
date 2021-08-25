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


//Test Handle
func (this *PingRouter) Handle(request yiface.IRequest) {
	fmt.Println("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	//回写数据
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}
//HelloYgoRouter Handle
type HelloYgoRouter struct {
	ynet.BaseRouter
}

func (this *HelloYgoRouter) Handle(request yiface.IRequest) {
	fmt.Println("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	fmt.Println("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Ygo Router V0.6"))
	if err != nil {
		fmt.Println(err)
	}
}

func main() {

	//1 创建一个server 句柄 s
	s := ynet.NewServer()

	s.AddRouter(0,&PingRouter{})
	s.AddRouter(1,&HelloYgoRouter{})
	//2 开启服务
	s.Server()
}
