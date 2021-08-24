package yiface

type IPack interface {
	GetHeadLen() uint32                //获取包头长度方法
	Pack(msg IMessage) ([]byte, error) //封包
	UnPack([]byte) (IMessage, error)   //拆包
}
