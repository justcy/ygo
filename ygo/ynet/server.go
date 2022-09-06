package ynet

import (
	"context"
	"errors"
	"fmt"
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
	Name      string
	IPVersion string
	IP        string
	Port      int
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler yiface.IMsgHandle
	//当前Server的链接管理器
	ConnMgr       yiface.IConnManager
	OnServerStart func(server yiface.IServer)
	OnServerStop  func(server yiface.IServer)
	//该Server的连接创建时Hook函数
	OnConnStart func(conn yiface.IConnection)
	//该Server的连接断开时的Hook函数
	OnConnStop func(conn yiface.IConnection)
	packet     yiface.IPack
	sigs     chan os.Signal
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

//============== 定义当前客户端链接的handle api ===========
func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务
	ylog.Debug("[Conn Handle] CallBackToClient ... ")
	if _, err := conn.Write(data[:cnt]); err != nil {
		ylog.Infof("write back buf err %s", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

func (s *Server) Start() {
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
	s.listenSignal(context.Background(), s)
}

func (s *Server) Server() {
	s.Start()
	//TODO Server.Serve() 是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加
	s.CallOnServerStart(s)

	//阻塞,否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}
func (s *Server) listenSignal(ctx context.Context, server yiface.IServer) {
	s.sigs = make(chan os.Signal, 1)
	signal.Notify(s.sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	select {
	case <-s.sigs:
		server.Stop()
		os.Exit(0)
	}
}

func (s *Server) Stop() {
	fmt.Printf("[STOP] %s exited!", s.Name)
	//将其他需要清理的连接信息或者其他信息 也要一并停止或者清理
	s.CallOnServerStop(s)
	s.ConnMgr.ClearConn()
}

func NewServer() yiface.IServer {
	utils.GlobalObject.Reload()
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), //创建ConnManager
		packet:     NewDataPack(),
	}
	return s
}
