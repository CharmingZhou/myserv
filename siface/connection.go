package siface

import "net"

// 定义连接接口
type Connection interface {
	Start()                         //启动连接，让当前连接开始工作
	Stop()                          //停止连接，结束当前连接状态M
	GetTCPConnection() *net.TCPConn //从当前连接获取原始的socket TCPConn
	GetConnID() uint32              //获取当前连接ID
	RemoteAddr() net.Addr           //获取远程客户端地址信息
	SendMsg(msgId uint32, data []byte) error
	SendBuffMsg(msgId uint32, data []byte) error //添加带缓冲发送消息接口

	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
}

// 定义一个统一处理连接业务的接口
type HandFunc func(*net.TCPConn, []byte, int) error
