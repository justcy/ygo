package ynet

import (
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
	"strconv"
)

type MsgHandle struct {
	Apis map[uint32] yiface.IRouter
}

func (mh *MsgHandle) DoMsgHandler(request yiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId = ", request.GetMsgId(), " is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.AfterHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router yiface.IRouter) {
	//1 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api , msgId = " + strconv.Itoa(int(msgId)))
	}
	//2 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: map[uint32]yiface.IRouter{},
	}
}

