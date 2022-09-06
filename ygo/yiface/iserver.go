package yiface

type IServer interface {
	//启动服务器
	Start()
	//开启业务服务代码
	Server()
	//停止服务器
	Stop()
	//路由功能:给当前服务注册一个路由业务方法，供客户端连接处理时使用
	AddRouter(msgId uint32,router IRouter)
	//得到链接管理
	GetConnMgr() IConnManager
	//设置该Server启动时Hook函数
	SetOnServerStart(func (IServer))
	SetOnServerStop(func (IServer))
	//调用连接OnServerStart Hook函数
	CallOnServerStart(s IServer)
	//调用连接OnServerStop Hook函数
	CallOnServerStop(s IServer)
	//设置该Server的连接创建时Hook函数
	SetOnConnStart(func (IConnection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func (IConnection))
	//调用连接OnConnStart Hook函数
	CallOnConnStart(conn IConnection)
	//调用连接OnConnStop Hook函数
	CallOnConnStop(conn IConnection)

	Packet() IPack
}