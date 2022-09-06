package yiface

type IRequest interface {
	//GetServer() IServer //获取所属服务器ID
	GetConnection() IConnection //获取请求链接
	GetData() []byte            //获取请求消息
	GetMsgId() uint32           //获取消息ID
}
