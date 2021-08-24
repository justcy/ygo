package yiface

import "net"

type IConnection interface {
	//启动连接
	Start()
	//停止连接
	Stop()
	//从当前连接获取Socket TCPConn
	GetTCPConnection() *net.TCPConn
	//获取当前连接的ID
	GetConnId() uint32
	//获取远程客户端地址信息
	RemoteAddr() net.Addr
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn,[]byte,int) error