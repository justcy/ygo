package yiface

type IRequest interface {
	GetConnection() IConnection //获取请求链接
	GetData() []byte            //获取请求消息
}
