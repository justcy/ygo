package ynet

import "github.com/justcy/ygo/ygo/yiface"

type Request struct {
	Conn yiface.IConnection //已经建立好的连接
	Msg  yiface.IMessage    //客户端请求的数据
}

func (r *Request) GetConnection() yiface.IConnection {

	return r.Conn
}

func (r *Request) GetData() []byte {
	return r.Msg.GetData()
}
func (r *Request) GetMsgId() uint32 {
	return r.Msg.GetMsgId()
}

