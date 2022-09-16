package main

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yclient"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"github.com/justcy/ygo/ygo/ynet"
	"time"
)

//ping test 自定义路由
type TestRouter struct {
	ynet.BaseRouter //一定要先基础BaseRouter
}

//Test Handle
func (this *TestRouter) Handle(request yiface.IRequest) {
	ylog.Info("Call TestRouter Handle")
	//先读取客户端的数据，再回写ping...ping...ping
	ylog.Infof("recv from server : msgId=%d, data=%s", request.GetMsgId(), string(request.GetData()))
	//回写数据
	//conn := request.GetConnection()
	err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	if err != nil {
		ylog.Error(err)
	}
}

func main() {
	ylog.SetLogPath("/Users/justcy/Documents/Develop/go/src/github.com/justcy/ygo/examples/log/client",ylog.LogSplitDay)
	ylog.Debug("start")
	client := yclient.NewClient(fmt.Sprintf("%s:%d","127.0.0.1",7777))
	client.AddRouter(2,&TestRouter{})
	client.Start()
	time.Sleep(5 * time.Second)
	ylog.Debug(client.GetConn())
	client.GetConn().SendMsg(1, []byte("这是一个测试"))
}
