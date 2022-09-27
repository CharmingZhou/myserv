package main

import (
	"fmt"
	"io"
	"net"
	"time"

	"github.com/CharmingZhou/myserv/siface"
	"github.com/CharmingZhou/myserv/snet"
)

type PingRouter struct {
	snet.BaseRouter
}

func (this *PingRouter) Handle(request siface.Request) {
	fmt.Println("Call PingRouter Hanle")
	//先读取客户端的数据，再写回ping...ping....ping
	fmt.Println("recv from client: msgid=", request.GetMsgId(), ", data=", string(request.GetData()))
	//写回数据
	err := request.GetConnection().SendMsg(0, []byte("ping...ping...ping"))
	if err != nil {
		fmt.Println(err)
	}
}

type HelloMyServRouter struct {
	snet.BaseRouter
}

func (this *HelloMyServRouter) Handle(request siface.Request) {
	fmt.Println("Call HelloMyServRouter Handle")
	fmt.Println("recv from client: msgId=", request.GetMsgId(), ",data=", string(request.GetData()))

	err := request.GetConnection().SendMsg(1, []byte("Hello MyServ Router V0.8"))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionBegin(conn siface.Connection) {
	fmt.Println("DoConnection is Called...")

	//设置两个连接属性，在连接创建之后
	fmt.Println("Set conn Name, Home done!")
	conn.SetProperty("Name", "Charming")
	conn.SetProperty("Home", "https://charmingz.com")

	err := conn.SendMsg(2, []byte("DoConnection BEGIN..."))
	if err != nil {
		fmt.Println(err)
	}
}

func DoConnectionLost(conn siface.Connection) {
	//在连接销毁之前，查询conn的Name，Home属性
	if name, err := conn.GetProperty("Name"); err == nil {
		fmt.Println("Conn Property Name=", name)
	}

	if home, err := conn.GetProperty("Home"); err == nil {
		fmt.Println("Conn Property Home =", home)
	}

	fmt.Println("DoConnectionLost is Called...")
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
	i := 0
	for {
		i++
		//发封包message消息
		dp := snet.NewDataPack()
		msgId := uint32(i % 2)
		msg, _ := dp.Pack(snet.NewMsgPackage(msgId, []byte("MyServ V0.5 Client Test Message")))
		_, err := conn.Write(msg)
		if err != nil {
			fmt.Println("write error err ", err)
			return
		}

		//先读出流中的head部分
		headData := make([]byte, dp.GetHeadLen())
		_, err = io.ReadFull(conn, headData) //ReadFull会把msg填充满为止
		if err != nil {
			fmt.Println("read head error")
			break
		}

		//将headData字节流 拆分到msg中
		msgHead, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("server unpack err:", err)
			return
		}

		if msgHead.GetDataLen() > 0 {
			//msg是有data数据的， 需要再次读取data数据
			msg := msgHead.(*snet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			//根据dataLen从io中读取字节流
			_, err := io.ReadFull(conn, msg.Data)
			if err != nil {
				fmt.Println("server unpack data err:", err)
				return
			}
			fmt.Println("===> Recv Msg: ID=", msg.Id, ", len=", msg.DataLen, ",data=", string(msg.Data))
		}
		time.Sleep(1 * time.Second)
	}
	time.Sleep(1 * time.Second)
}

func main() {
	s := snet.NewServer("[MyServ V0.9]")

	//注册连接hook回调函数
	s.SetOnConnStart(DoConnectionBegin)
	s.SetOnConnStop(DoConnectionLost)

	//配置路由
	s.AddRouter(0, &PingRouter{})
	s.AddRouter(1, &HelloMyServRouter{})

	go func() {
		clientTest()
	}()
	// 2.开启服务
	s.Serve()
}
