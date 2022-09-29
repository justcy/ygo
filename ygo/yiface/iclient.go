package yiface

type IClient interface {
	//启动服务器
	Start()
	//停止服务器
	Stop()
	//路由功能:给当前服务注册一个路由业务方法，供客户端连接处理时使用
	AddRouter(msgId uint32,router IRouter)
	GetConn() IConnection
	TickAck() bool
	SetAckMsg(heart []byte)
	GetActMsg() []byte
}
