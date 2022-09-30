package yiface

import (
	"net"
)

type IClient interface {
	Start()
	Stop(bool)
	GetConnection() *net.TCPConn
	AddRouter(msgId uint32, router IRouter)
	Send(msgId uint32,data []byte) error
	GetProperty(string) (interface{}, error)
	SetProperty(string, interface{})
	RemoveProperty(string)
	TickAck() bool
}
