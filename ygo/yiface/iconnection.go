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
	//直接将Message消息发送给远程TCP客户端
	SendMsg(msgId uint32,data []byte) error
	//直接将Message数据发送给远程的TCP客户端(有缓冲)
	SendBuffMsg(msgId uint32, data []byte) error   //添加带缓冲发送消息接口

	//设置链接属性
	SetProperty(key string, value interface{})
	//获取链接属性
	GetProperty(key string)(interface{}, error)
	//移除链接属性
	RemoveProperty(key string)
}

//定义一个统一处理链接业务的接口
type HandFunc func(*net.TCPConn,[]byte,int) error