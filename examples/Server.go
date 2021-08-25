package main

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ynet"
)

//ping test 自定义路由
type PingRouter struct {
	ynet.BaseRouter //一定要先基础BaseRouter
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

//创建连接的时候执行
func DoConnectionBegin(conn yiface.IConnection) {
	fmt.Println("DoConnecionBegin is Called ... ")

	//=============设置两个链接属性，在连接创建之后===========
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Justcy")
	conn.SetProperty("Home", "http://blog.kanter.cn")
	//===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn yiface.IConnection) {
	if name, err:= conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name = ", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home = ", home)
	}

	fmt.Println("DoConneciotnLost is Called ... ")
}

func main() {

	//1 创建一个server 句柄 s
	s := ynet.NewServer()

	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloYgoRouter{})
	//2 开启服务
	s.Server()
}
