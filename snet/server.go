package snet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/CharmingZhou/myserv/utils"

	"github.com/CharmingZhou/myserv/siface"
)

type Server struct {
	Name      string //服务器的名称
	IPVersion string //tcp4 or other
	IP        string //服务绑定的IP地址
	Port      int    //服务绑定的端口
	//Router    siface.Router //当前Server由用户绑定的回调router，也就是Server注册的链接对应的处理业务
	msgHandler siface.MsgHandler //当前Server的消息管理模块，用来绑定MsgId和对应的处理方法
	ConnMgr    siface.ConnManager

	OnConnStart func(conn siface.Connection)
	OnConnStop  func(conn siface.Connection)
}

func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
	//回显业务

	fmt.Println("[Conn Handle] CallbackToClient...")
	if _, err := conn.Write(data[:cnt]); err != nil {
		fmt.Println("write back buf err ", err)
		return errors.New("CallBackToClient error")
	}
	return nil
}

// 开启网络服务
func (s *Server) Start() {
	fmt.Printf("[START] Server listener at IP:%s, Port:%d, is starting\n", s.IP, s.Port)
	fmt.Printf("[MyServ] Version:%s, MaxConn:%d, MaxpacketSize:%d\n",
		utils.GlobalObject.Version,
		utils.GlobalObject.MaxConn,
		utils.GlobalObject.MaxPacketSize)
	//开启一个go去做服务端Listener业务
	go func() {
		//0 启动worker工作池机制
		s.msgHandler.StartWorkerPool()

		//1 获取一个TCP的Addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr err:", err)
			return
		}
		//2 监听服务器地址
		listenner, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, "err", err)
			return
		}
		//已经监听成功
		fmt.Println("start myserv server ", s.Name, "succ, now listenning...")

		//TODO:serer.go应该有一个自动生成ID的方法
		var cid uint32
		cid = 0

		//3 启动server网络连接业务
		for {
			//3.1 阻塞等待客户端建立连接请求
			conn, err := listenner.AcceptTCP()
			if err != nil {
				fmt.Println("Accept err", err)
				continue
			}

			//3.2 TODO: Server.Start()设置服务器最大连接控制，如果超过最大连接，那么则关闭此新的连接
			if s.ConnMgr.Len() >= utils.GlobalObject.MaxConn {
				conn.Close()
				continue
			}

			//3.3 处理该新连接请求的业务方法，此时应该有handler和conn的绑定的
			dealConn := NewConnection(s, conn, cid, s.msgHandler)
			cid++

			//3.4 启动当前连接的处理业务
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] mysvr server ,name ", s.Name)
	//TODO Server.Stop()将其他需要清理的连接信息或者其他信息，也要一并停止或者清理
	s.ConnMgr.ClearConn()
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve()是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞，否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}
func (s *Server) AddRouter(msgId uint32, router siface.Router) {
	s.msgHandler.AddRouter(msgId, router)
	fmt.Println("Add Router succ! ")
}

// 创建一个服务器句柄
func NewServer(name string) siface.Server {
	//先初始化全局配置文件
	utils.GlobalObject.Reload()

	s := &Server{
		Name:      utils.GlobalObject.Name,
		IPVersion: "tcp4",
		IP:        utils.GlobalObject.Host,
		Port:      utils.GlobalObject.TcpPort,
		//Router:    nil,
		msgHandler: NewMsgHandle(),
		ConnMgr:    NewConnManager(), //创建ConnManager
	}
	return s
}
func (s *Server) GetConnMgr() siface.ConnManager {
	return s.ConnMgr
}

func (s *Server) SetOnConnStart(hookFunc func(connection siface.Connection)) {
	s.OnConnStart = hookFunc
}

func (s *Server) SetOnConnStop(hookFunc func(connection siface.Connection)) {
	s.OnConnStop = hookFunc
}

func (s *Server) CallOnConnStart(conn siface.Connection) {
	if s.OnConnStart != nil {
		fmt.Println("--->CallOnConnStart...")
		s.OnConnStart(conn)
	}
}

func (s *Server) CallOnConnStop(conn siface.Connection) {
	if s.OnConnStop != nil {
		fmt.Println("----> CallOnConnStop....")
		s.OnConnStop(conn)
	}
}
