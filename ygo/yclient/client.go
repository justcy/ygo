package yclient

import (
	"context"
	"github.com/hashicorp/go-uuid"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"github.com/justcy/ygo/ygo/ynet"
	"net"
	"sync"
	"time"
)

type client struct {
	Id   string
	Addr string
	Conn yiface.IConnection
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler yiface.IMsgHandle
	sync.RWMutex
	//告知该链接已经退出/停止的channel
	ctx      context.Context
	cancel   context.CancelFunc
	isClosed bool

	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	wg              sync.WaitGroup
}

var err error

func (c *client) dial() net.Conn {
	for {
		conn, err := net.Dial("tcp", c.Addr)
		if err == nil || c.isClosed {
			return conn
		}
		time.Sleep(c.ConnectInterval)
		continue
	}
}

func (c *client) Start() {
	ylog.Debug(c.Addr)
	c.wg.Add(1)
	go c.connect()
}

func (c *client) Stop() {
	//1. 如果当前链接已经关闭
	if c.isClosed == true {
		return
	}
	c.isClosed = true
	c.GetConn().Stop()
	c.wg.Wait()

}

func (c *client) AddRouter(msgId uint32, router yiface.IRouter) {
	c.msgHandler.AddRouter(msgId, router)
}

func (c *client) GetConn() yiface.IConnection {
	return c.Conn
}

func (c *client) connect() {
	defer c.wg.Done()
reconnect:
	conn := c.dial()
	if conn == nil {
		return
	}
	c.Lock()
	if c.isClosed {
		c.Unlock()
		conn.Close()
		return
	}
	c.Unlock()
	c.Conn = NewConnection(conn.(*net.TCPConn), 1, c.msgHandler)
	c.Conn.Start()
	if c.AutoReconnect {
		time.Sleep(c.ConnectInterval)
		goto reconnect
	}
}

func NewClient(address string) *client {
	uuid, _ := uuid.GenerateUUID()
	if err != nil {
		ylog.Errorf("Generate Server UUID %s", err)
	}
	c := &client{
		Id:              uuid,
		Addr:            address,
		isClosed:        false,
		msgHandler:      ynet.NewMsgHandle(),
		ConnectInterval: 5 * time.Second,
		PendingWriteNum: 2048,
		AutoReconnect:   true,
	}
	return c
}
