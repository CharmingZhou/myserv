package snet

import (
	"fmt"
	"strconv"

	"github.com/CharmingZhou/myserv/siface"
)

type MsgHandle struct {
	Apis map[uint32]siface.Router //存放每个MsgId所对应的处理方法的map属性
}

func NewMsgHandle() *MsgHandle {
	return &MsgHandle{
		Apis: make(map[uint32]siface.Router),
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
