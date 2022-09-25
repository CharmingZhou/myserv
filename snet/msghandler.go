package snet

import (
	"fmt"
	"strconv"

	"github.com/CharmingZhou/myserv/utils"

	"github.com/CharmingZhou/myserv/siface"
)

type MsgHandle struct {
	Apis           map[uint32]siface.Router //存放每个MsgId所对应的处理方法的map属性
	WorkerPoolSize uint32                   //业务工作Worker池的数量
	TaskQueue      []chan siface.Request    //Work负责取消任务的消息队列
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis:           make(map[uint32]siface.Router),
		WorkerPoolSize: utils.GlobalObject.WorkPoolSize,
		//一个worker对应一个queue
		TaskQueue: make([]chan siface.Request, utils.GlobalObject.WorkPoolSize),
	}
}

func (mh *MsgHandle) DoMsgHandler(request siface.Request) {
	handler, ok := mh.Apis[request.GetMsgId()]
	if !ok {
		fmt.Println("api msgId=", request.GetMsgId(), "is not FOUND!")
		return
	}
	//执行对应处理方法
	handler.PreHandle(request)
	handler.Handle(request)
	handler.PostHandle(request)
}

func (mh *MsgHandle) AddRouter(msgId uint32, router siface.Router) {
	//1. 判断当前msg绑定的API处理方法是否已经存在
	if _, ok := mh.Apis[msgId]; ok {
		panic("repeated api, msgId = " + strconv.Itoa(int(msgId)))
	}
	//2. 添加msg与api的绑定关系
	mh.Apis[msgId] = router
	fmt.Println("Add api msgId = ", msgId)
}

func (mh *MsgHandle) StartOneWorker(workerID int, taskQueue chan siface.Request) {
	fmt.Println("Worker ID=", workerID, " is started.")
	//不断的等待队列中的消息
	for {
		select {
		//有消息则取出队列的Request, 并执行绑定的业务方法
		case request := <-taskQueue:
			mh.DoMsgHandler(request)
		}
	}
}

// 启动worker工作池
func (mh *MsgHandle) StartWorkerPool() {
	//遍历需要启动worker的数量，依次启动
	for i := 0; i < int(mh.WorkerPoolSize); i++ {
		//一个worker被启动
		//给当前worker对应的任务队列开辟空间
		mh.TaskQueue[i] = make(chan siface.Request, utils.GlobalObject.MaxWorkerTaskLen)
		//启动当前Worker，阻塞的等待对应的任务队列是否有消息传递进来
		go mh.StartOneWorker(i, mh.TaskQueue[i])
	}
}

// 将消息交给TaskQueue，由worker进行处理
func (mh *MsgHandle) SendMsgToTaskQueue(request siface.Request) {
	//根据ConnID来分配当前的连接应该由哪个worker负责处理
	//轮训的平均分配法则

	//得到需要处理此条连接的workerID
	workerID := request.GetConnection().GetConnID() % mh.WorkerPoolSize
	fmt.Println("Add ConnID=", request.GetConnection().GetConnID(), " request msgID=", request.GetMsgId(),
		"to workerID=", workerID)
	//将请求消息发送给任务队列
	mh.TaskQueue[workerID] <- request
}
