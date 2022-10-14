package ynet

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-uuid"
	"github.com/justcy/ygo/ygo/registry"
	"github.com/justcy/ygo/ygo/registry/iface"
	"github.com/justcy/ygo/ygo/utils"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	Id        string
	Name      string
	IPVersion string
	IP        string
	Port      int16
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler yiface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr yiface.IConnManager
	//该Server启动时Hook函数
	OnServerStart func(server yiface.IServer)
	//该Server退出时Hook函数
	OnServerStop func(server yiface.IServer)
	//该Server的连接创建时Hook函数
	OnConnStart func(conn yiface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn yiface.IConnection)

	//tick函数hook
	OnTick func(tick time.Time)

	packet yiface.IPack
	//接收关闭信号
	sigs chan os.Signal

	ctx    context.Context
	cancel context.CancelFunc

	//
	tickChan chan bool

	Tick100MSec   int64 //100毫秒
	tick300MSec   int64 //300毫秒
	tickOneSec    int64 //1秒
	tickFiveSec   int64 //5秒
	tickThirtySec int64 //30秒
	tickSixtySec  int64 //60秒
	tickFiveMin   int64 //5分钟
	//Client chan yiface.IClient
	Client        map[string]yiface.IClient
}

func (s *Server) Packet() yiface.IPack {
	return s.packet
}

func (s *Server) SetOnConnStart(hookFunc func(yiface.IConnection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(yiface.IConnection)) {
	s.OnConnStop = hookFunc
}
func (s *Server) SetOnServerStart(hookFunc func(yiface.IServer)) {
	s.OnServerStart = hookFunc
}

func (s *Server) SetOnServerStop(hookFunc func(yiface.IServer)) {
	s.OnServerStop = hookFunc
}
func (s *Server) SetOnTick(hookFunc func(time.Time)) {
	s.OnTick = hookFunc
}

func (s *Server) CallOnConnStart(conn yiface.IConnection) {
	if s.OnConnStart != nil {
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn yiface.IConnection) {
	if s.OnConnStop != nil {
		s.OnConnStop(conn)
	}
}

func (s *Server) CallOnServerStart(server yiface.IServer) {
	if s.OnServerStart != nil {
		s.OnServerStart(server)
	}
}

func (s *Server) CallOnServerStop(server yiface.IServer) {
	if s.OnServerStop != nil {
		s.OnServerStop(server)
	}
}

func (s *Server) GetConnMgr() yiface.IConnManager {
	return s.ConnMgr
}

func (s *Server) AddRouter(msgId uint32, router yiface.IRouter) {
	ylog.Debug("Add Router success !")
	s.msgHandler.AddRouter(msgId, router)
}

func (s *Server) AddClient(key string, c yiface.IClient) {
	s.Client[key] = c
	s.Client[key].Start()
}
func (s *Server) GetClient(key string)yiface.IClient  {
	if 	s.Client[key] == nil{
		ylog.Errorf("client not find %s",key)
	}
	return s.Client[key]
}
func (s *Server) Start() {
	s.ctx, s.cancel = context.WithCancel(context.Background())
	fmt.Printf("[START] Server listenner at IP :%s,Port %d,is Starting\n", s.IP, s.Port)
	fmt.Printf("[Ygo] Version:%s,MaxConn:%d,MaxPacketSize:%d\n", utils.GlobalObject.Version, utils.GlobalObject.MaxConn, utils.GlobalObject.MaxPacketSize)
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()
		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			ylog.Errorf("resolve tcp addr err: %s", err)
			return
		}

		//2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			ylog.Errorf("listen %s, err %s", s.IPVersion, err)
			return
		}

		//已经监听成功
		fmt.Println("Start ", s.Name, " succ, now listening...")

		//TODO server.go 应该有一个自动生成ID的方法
		var cid uint32
		cid = 0
		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				ylog.Errorf("Accept err %s", err)
				continue
			}
			//3.2 设置服务器最大连接控制,如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}
			//3.3 处理该新连接请求的 业务 方法， 此时应该有 handler 和 conn是绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++
			go dealConn.Start()
		}

	}()
}

func (s *Server) Server() {
	s.sigs = make(chan os.Signal, 1)
	signal.Notify(s.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	s.CallOnServerStart(s)

	if utils.GlobalObject.Tick {
		s.tickChan = make(chan bool)
		s.Tick100MSec = time.Now().UnixNano() / 1e6
		s.tick300MSec = time.Now().UnixNano() / 1e6
		s.tickOneSec = time.Now().Unix()
		s.tickFiveSec = s.tickOneSec + 5
		s.tickThirtySec = s.tickOneSec + 30
		s.tickSixtySec = s.tickOneSec + 60
		s.tickFiveMin = s.tickOneSec + 300

		go func() {
			mtick := time.NewTicker(100 * time.Millisecond)
			for _ = range mtick.C {
				s.tickChan <- true
			}
		}()
	}

	if utils.GlobalObject.ConsulAddress != "" {
		var err error
		s.Id, err = uuid.GenerateUUID()
		if err != nil {
			ylog.Errorf("Generate Server UUID %s", err)
		}
		consulRegister := &registry.ConsulRegistry{}
		consulRegister.Register(iface.Service{
			Id:       s.Id,
			Name:     s.Name,
			Version:  "v0.0.1",
			//Address:  "115.28.133.188",
			Address:  "127.0.0.1",
			Port:     s.Port,
			Metadata: map[string]string{"tags": "111", "role": "gate"},
		})
	}
	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		select {
		case <-s.tickChan:
			s.Tick(time.Now())
		case <-s.sigs:
			s.Stop()
			os.Exit(0)
		}
	}
}

//func (s *Server) listenSignal(ctx context.Context, server yiface.IServer) {
//
//	select {
//	case <-s.sigs:
//		server.Stop()
//		os.Exit(0)
//	}
//}

func (s *Server) Stop() {
	fmt.Printf("[STOP] %s exited!", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.CallOnServerStop(s)
	s.ConnMgr.ClearConn()
	consulRegister := &registry.ConsulRegistry{}
	consulRegister.UnRegisterById(s.Id)

	for _, client := range s.Client {
		ylog.Debugf("stop client %v",client)
		client.Stop(true)
	}
	s.cancel()
}

func (s *Server) Tick(tick time.Time) {
	mSec := tick.UnixNano() / 1e6
	if mSec >= s.Tick100MSec { //100ms
		s.Tick100MSec = mSec + 100
	}
	if mSec >= s.tick300MSec { //300ms
		s.tick300MSec = mSec + 300
	}
	nNow := tick.Unix()
	if nNow >= s.tickOneSec {
		s.tickOneSec = nNow + 1
	}
	if nNow >= s.tickFiveSec { //5秒
		s.sendClientAck()
		s.tickFiveSec = nNow + 5
	}
	if nNow >= s.tickThirtySec { //30秒
		s.tickThirtySec = nNow + 30
	}
	if nNow >= s.tickSixtySec { //60秒
		s.tickSixtySec = nNow + 60
	}
	if nNow >= s.tickFiveMin { //300秒
		s.ConnMgr.Tick()
		s.tickFiveMin = nNow + 300
	}
	if s.OnTick != nil {
		s.OnTick(tick)
	}
}

func (s *Server) sendClientAck() {
	if s.Client == nil {
		return
	}
	for _, client := range s.Client {
		if !client.TickAck(){
			continue
		}
		if heart, err := client.GetProperty(HEART_MSG); err == nil {
			client.GetConnection().Write(heart.([]byte))
		}

	}
}

func (s *Server) GetInfo() *Server {
	return s
}
func (s *Server) GetCtx() context.Context {
	return s.ctx
}

func NewServer(conf string) yiface.IServer {
	utils.GlobalObject.Reload(conf)
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), //创建ConnManager
		packet:     NewDataPack(),
		Client: make(map[string] yiface.IClient,1),
	}
	return s
}
