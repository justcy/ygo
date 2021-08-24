package ynet

import "github.com/justcy/ygo/ygo/yiface"

type Request struct {
	conn yiface.IConnection//已经建立好的连接
	msg yiface.IMessage//客户端请求的数据
}

func (r *Request) GetConnection() yiface.IConnection {

	return r.conn
}

func (r *Request) GetData() []byte {
	return r.msg.GetData()
}
func (r *Request) GetMsgId() uint32 {
	return r.msg.GetMsgId()
}

