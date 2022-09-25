package main

import (
	"fmt"
	"net"
	"time"

	"github.com/CharmingZhou/myserv/siface"
	"github.com/CharmingZhou/myserv/snet"
)

type PingRouter struct {
	snet.BaseRouter
}

func (this *PingRouter) PreHandle(request siface.Request) {
	fmt.Println("Call Router PreHandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("before ping...\n"))
	if err != nil {
		fmt.Println(fmt.Println("call back ping error"))
	}
}

func (this *PingRouter) Handle(request siface.Request) {
	fmt.Println("Call PingRouter Hanle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("ping...ping...ping\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func (this *PingRouter) PostHandle(request siface.Request) {
	fmt.Println("Call Router Posthandle")
	_, err := request.GetConnection().GetTCPConnection().Write([]byte("After ping...\n"))
	if err != nil {
		fmt.Println("call back ping ping ping error")
	}
}

func clientTest() {
	fmt.Println("Client Test...start")
	//3秒之后发起测试请求，给服务端开启服务的机会
	time.Sleep(3 * time.Second)

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client start err, exit!")
		return
	}
	for {
		_, err := conn.Write([]byte("myserv V0.3"))
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buf err")
			return
		}
		fmt.Printf("server call back:%s, cnt=%d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}
}

func main() {
	s := snet.NewServer("[myserv V0.3]")
	s.AddRouter(&PingRouter{})
	go func() {
		clientTest()
	}()
	// 2.开启服务
	s.Serve()
}
