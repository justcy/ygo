package yclient

import (
	"errors"
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"github.com/justcy/ygo/ygo/ynet"
	"net"
	"sync"
	"time"
)

type client struct {
	Id   string
	Addr *net.TCPAddr
	Conn yiface.IConnection
	//当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	msgHandler yiface.IMsgHandle
	sync.RWMutex
	//告知该链接已经退出/停止的channel
	isClosed bool

	retryInterval   int
	ConnectInterval time.Duration
	PendingWriteNum int
	AutoReconnect   bool
	isTickAck       bool
	ActMsg          []byte

	//链接属性
	property map[string]interface{}
	//保护链接属性修改的锁
	propertyLock sync.RWMutex
}

var (
	HEART_MSG = "heart_msg"
)

func NewClient(id string, ip string, port int) *client {
	addr := &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
		Zone: "",
	}
	c := &client{
		Id:              id,
		Addr:            addr,
		isClosed:        false,
		msgHandler:      ynet.NewMsgHandle(),
		retryInterval:   1024,
		ConnectInterval: 5 * time.Second,
		PendingWriteNum: 2048,
		property:        make(map[string]interface{}, 1),
		AutoReconnect:   true,
		isTickAck:       true,
		ActMsg:          nil,
	}
	return c
}
func (c *client) Start() {
	if c.Conn == nil {
		c.connection()
	}
	go c.Conn.Start()
}

func (c *client) Stop(b bool) {

}
func (c *client) TickAck() bool {
	return c.isTickAck
}
func (c *client) GetConnection() *net.TCPConn {
	return c.Conn.GetTCPConnection()
}

func (c *client) AddRouter(msgId uint32, router yiface.IRouter) {
	c.RLock()
	defer c.RUnlock()
	c.msgHandler.AddRouter(msgId, router)
}

func (c *client) Send(msgId uint32, data []byte) error {
	c.RLock()
	defer c.RUnlock()

	return c.Conn.SendMsg(msgId, data)
}

func (c *client) GetProperty(s string) (interface{}, error) {
	c.propertyLock.RLock()
	defer c.propertyLock.RUnlock()

	value, ok := c.property[s]
	if ok {
		return value, nil
	} else {
		return nil, errors.New("no property in connection")
	}
}

func (c *client) SetProperty(s string, i interface{}) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	c.property[s] = i
}

func (c *client) RemoveProperty(s string) {
	c.propertyLock.Lock()
	defer c.propertyLock.Unlock()
	delete(c.property, s)
}

func (c *client) connection() {
	c.Lock()
	defer c.Unlock()
	if c.Conn != nil {
		return
	}
	for i := 1; i < c.retryInterval; i++ {
		ylog.Info("retry time ", i)
		conn, err := net.DialTCP("tcp", nil, c.Addr)
		if err == nil {
			c.Conn = NewConnection(c, conn, c.msgHandler)
			return
		} else {
			d, err := time.ParseDuration(fmt.Sprintf("%ds", c.retryInterval))
			if err != nil {
				time.Sleep(c.ConnectInterval)
			} else {
				time.Sleep(d)
			}
		}
	}

}
