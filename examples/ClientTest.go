package main

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yclient"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"github.com/justcy/ygo/ygo/ynet"
	"github.com/justcy/ygo/ygo/ytimer"
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
	//err := request.GetConnection().SendMsg(1, []byte("ping...ping...ping"))
	//if err != nil {
	//	ylog.Error(err)
	//}
}
func SayHello(t *ytimer.Task) {
	t.Args[0].(*ynet.Connection).SendMsg(1,[]byte(t.Args[1].(string)))
	fmt.Printf("%s 执行次数 %d hello %s,%s\n", time.Now(), t.Args[0], t.Args[1])
}
func main() {
	ylog.SetLogPath("/Users/justcy/Documents/Develop/go/src/github.com/justcy/ygo/examples/log/client",ylog.LogSplitDay)
	ylog.Debug("start")
	client := yclient.NewClient(fmt.Sprintf("%s:%d","127.0.0.1",7777))
	client.AddRouter(2,&TestRouter{})
	client.Start()
	time.Sleep(5 * time.Second)
	ylog.Debug(client.GetConn())

	//tw, _ := ytimer.NewTimeWheel(1*time.Second, 10, ytimer.TickSafeMode())
	//tw.AddCron(2*time.Second, ytimer.ModeIsAsync, SayHello, []interface{}{client.GetConn().GetTCPConnection(),"message type"})
	//tw.Start()

	for {
		client.GetConn().SendMsg(1, []byte("this is a Test！！"))
		time.Sleep(5 * time.Second)
	}

}
