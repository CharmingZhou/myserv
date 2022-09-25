package siface

type MsgHandler interface {
	DoMsgHandler(request Request)
	AddRouter(msgId uint32, router Router)
}
