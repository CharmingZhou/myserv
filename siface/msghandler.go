package siface

type MsgHandler interface {
	DoMsgHandler(request Request)
	AddRouter(msgId uint32, router Router)
	StartWorkerPool()                   //启动worker工作池
	SendMsgToTaskQueue(request Request) //将消息交给TaskQueue，由woker进行处理
}
