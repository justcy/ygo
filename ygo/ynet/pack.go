package ynet

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/justcy/ygo/ygo/utils"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
)

type Pack struct {
}

func (dp *Pack) GetHeadLen() uint32 {
	//Id uint32(4字节) + DataLen uint32 4字节
	return 8
}

func (dp *Pack) Pack(msg yiface.IMessage) ([]byte, error) {
	//创建一个存放bytes字节的缓冲
	dataBuff := bytes.NewBuffer([]byte{})
	//写datalen
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	//写msgId
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	//写data数据
	if err := binary.Write(dataBuff, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuff.Bytes(),nil
}

func (dp *Pack) UnPack(binaryData []byte) (yiface.IMessage, error) {
	//创建一个从输入二进制数据的ioReader
	dataBuff := bytes.NewReader(binaryData)

	//只解压head的信息，得到dataLen和msgID
	msg := &Message{}

	//读dataLen
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	//读msgID
	if err := binary.Read(dataBuff, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	//判断dataLen的长度是否超出我们允许的最大包长度

	if utils.GlobalObject.MaxPacketSize > 0 && msg.DataLen > utils.GlobalObject.MaxPacketSize  {
		ylog.Errorf("max:%d,Request:%d",utils.GlobalObject.MaxPacketSize,msg.DataLen)
		return nil, errors.New("Too large msg data recieved")
	}
	//这里只需要把head的数据拆包出来就可以了，然后再通过head的长度，再从conn读取一次数据
	return msg, nil

}

func NewDataPack() *Pack {
	return &Pack{}
}
