package siface

type Server interface {
	Start()                                //启动服务器方法
	Stop()                                 //停止服务器方法
	Serve()                                //开启服务方法
	AddRouter(msgId uint32, router Router) //路由功能：给当前服务注册一个路由业务方法，供客户端处理使用
	//得到连接管理
	GetConnMgr() ConnManager
	//设置该Server的连接创建时的Hook函数
	SetOnConnStart(func(Connection))
	//设置该Server的连接断开时的Hook函数
	SetOnConnStop(func(Connection))
	CallOnConnStart(conn Connection)
	CallOnConnStop(conn Connection)
}
