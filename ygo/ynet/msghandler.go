package ynet

import (
	"github.com/justcy/ygo/ygo/utils"
	"github.com/justcy/ygo/ygo/yiface"
	"github.com/justcy/ygo/ygo/ylog"
	"strconv"
)

type MsgHandle struct {
	Apis           map[uint32]yiface.IRouter
	WorkerPoolSize uint32                 //业务工作Worker池的数量
	TaskQueue      []chan yiface.IRequest //Worker负责取任务的消息队列
}

//启动一个Worker工作流程
func (mh *MsgHandle) startOneWorker(workerId int, taskQueue chan yiface.IRequest) {
	ylog.Infof("Worker ID = %d is started.", workerId)
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request，并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依此启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan yiface.IRequest, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.startOneWorker(i, mh.TaskQueue[i])
	}
}

//将消息交给TaskQueue,由worker进行处理 @todo 实现其他算法
func (mh *MsgHandle) SendMsgToTaskQueue(request yiface.IRequest) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮询的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnId() % mh.WorkerPoolSize
	ylog.Infof("Add ConnId=%d,request msgId=%d,to workerId=%d", request.GetConnection().GetConnId(), request.GetMsgId(), workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}

func (mh *MsgHandle) DoMsgHandler(request yiface.IRequest) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		ylog.Errorf("api msgId = %d is not FOUND!", request.GetMsgId() )
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

	ylog.Infof("Add api msgId = %d", msgId)
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           map[uint32]yiface.IRouter{},
		WorkerPoolSize: utils.GlobalObject.WorkerPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan yiface.IRequest, utils.GlobalObject.WorkerPoolSize),
	}
}
