package main

import (
	"fmt"
	consul "github.com/hashicorp/consul/api"
	"github.com/justcy/ygo/ygo/registry"
	"github.com/justcy/ygo/ygo/yclient"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"github.com/justcy/ygo/ygo/ynet"
	"time"
)

//ping test 自定义路由
type PingRouter struct {
	ynet.BaseRouter //一定要先基础BaseRouter
}

//Test Handle
func (this *PingRouter) Handle(request yiface.IRequest) {
	ylog.Info("Call PingRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	ylog.Infof("recv from client : msgId=%d, data=%s", request.GetMsgId(), string(request.GetData()))
	//回写数据
	//conn := request.GetConnection()
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		ylog.Error(err)
	}
}

//HelloYgoRouter Handle
type HelloYgoRouter struct {
	ynet.BaseRouter
}

func (this *HelloYgoRouter) Handle(request yiface.IRequest) {
	ylog.Debug("Call HelloZinxRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	ylog.Debugf("recv from client : msgId=", request.GetMsgId(), ", data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello Ygo Router V0.6"))
	if err != nil {
		ylog.Error(err)
	}
}

//创建连接的时候执行
func DoConnectionBegin(conn yiface.IConnection) {
	ylog.Debug("DoConnecionBegin is Called ... ")
	//=============设置两个链接属性，在连接创建之后===========
	ylog.Info("Set conn Name, Home done!")
	conn.SetProperty("Name", "Justcy")
	conn.SetProperty("Home", "http://blog.kanter.cn")
	//===================================================

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		ylog.Error(err)
	}
}

//连接断开的时候执行
func DoConnectionLost(conn yiface.IConnection) {
	if name, err := conn.GetProperty("Name"); err == nil {
		ylog.Infof("Conn Property Name = %s", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		ylog.Infof("Conn Property Home = %s", home)
	}
	ylog.Debugf("DoConneciotnLost is Called ... ")
}
func ServerStart(server yiface.IServer) {
	consulRegister := &registry.ConsulRegistry{
		QueryOptions: &consul.QueryOptions{
			AllowStale: true,
		},
	}
	resp, _ := consulRegister.GetService("Demo server")
	ylog.Debugf("得到的服务列表 %d", len(resp))
	for _, service := range resp {
		ylog.Debugf("服务列表  %s,%s:%d", service.Name, service.Address, service.Port)
		ylog.Debugf("服务列表  %v", service)
		client := yclient.NewClient(fmt.Sprintf("%s:%d",service.Address,service.Port))
		client.AddRouter(2,&PingRouter{})
		//client.Start()
		//client.GetConn().SendMsg(1, []byte("这是一个测试"))
		//
		//client.Stop()
		key := service.Address + ":" + string(service.Port)

		server.AddClient(key, client)
		server.GetClient(key).GetConn().SendMsg(1, []byte("这是一个测试"))

	}
	ylog.Debugf("ServerStart is Called ... ")
}
func ServerStop(server yiface.IServer) {
	ylog.Debugf("ServerStop is Called ... ")
}

func MyTick(tick time.Time) {
	ylog.Debugf("MyTick called %d", tick)
}

func main() {
	//1 创建一个server 句柄 s
	s := ynet.NewServer()
	ylog.SetLogFile("./log", "666", ylog.LogSplitDay)
	//ylog.CloseDebug()
	s.SetOnServerStart(ServerStart)
	s.SetOnServerStop(ServerStop)
	//注册链接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)
	s.SetOnTick(MyTick)
	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloYgoRouter{})
	//2 开启服务
	s.Server()
}
