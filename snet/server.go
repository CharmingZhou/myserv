package snet

import (
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/CharmingZhou/myserv/siface"
)

type Server struct {
	Name      string
	IPVersion string
	IP        string
	Port      int
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
	//开启一个go去做服务端Listener业务
	go func() {
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

			//3.3 处理该新连接请求的业务方法，此时应该有handler和conn的绑定的
			dealConn := NewConnection(conn, cid, CallBackToClient)
			cid++

			//3.4 启动当前连接的处理业务
			go dealConn.Start()

		}
	}()
}

func (s *Server) Stop() {
	fmt.Println("[STOP] mysvr server ,name ", s.Name)
	//TODO Server.Stop()将其他需要清理的连接信息或者其他信息，也要一并停止或者清理
}

func (s *Server) Serve() {
	s.Start()
	//TODO Server.Serve()是否在启动服务的时候 还要处理其他的事情呢 可以在这里添加

	//阻塞，否则主Go退出， listenner的go将会退出
	for {
		time.Sleep(10 * time.Second)
	}
}

// 创建一个服务器句柄
func NewServer(name string) siface.Server {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      7777,
	}
	return s
}
