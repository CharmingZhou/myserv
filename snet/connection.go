package snet

import (
	"errors"
	"fmt"
	"io"
	"net"

	"github.com/CharmingZhou/myserv/siface"
)

type Connection struct {
	Conn     *net.TCPConn //当前连接的socket TCP套接字
	ConnID   uint32       //当前连接的ID 也可以作为SessionID，ID全局唯一
	isClosed bool         //当前连接的关闭状态

	handleAPI siface.HandFunc //该连接的处理方法api

	Router siface.Router //该连接的处理方法router

	ExitBuffChan chan bool //告知该连接已经退出/停止的channel
}

// 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, router siface.Router) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		Router:       router,
		ExitBuffChan: make(chan bool, 1),
	}
	return c
}

// 处理conn读数据的Goroutine
func (c *Connection) StartReader() {
	fmt.Println("Reader Goroutine is running")
	defer fmt.Println(c.RemoteAddr().String(), " conn reader exit!")
	defer c.Stop()

	for {
		//创建拆包解包的对象
		dp := NewDataPack()

		//读取我们最大的数据到buf中
		headData := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(c.GetTCPConnection(), headData); err != nil {
			fmt.Println("recv msg head err", err)
			c.ExitBuffChan <- true
			continue
		}

		//拆包，得到 msgid和datalen放在msg中
		msg, err := dp.Unpack(headData)
		if err != nil {
			fmt.Println("unpack error", err)
			c.ExitBuffChan <- true
			continue
		}

		//根据dataLen读取data，放在msg.Data中
		var data []byte
		if msg.GetDataLen() > 0 {
			data = make([]byte, msg.GetDataLen())
			if _, err := io.ReadFull(c.GetTCPConnection(), data); err != nil {
				fmt.Println("read msg data error ", err)
				c.ExitBuffChan <- true
				continue
			}
		}
		msg.SetData(data)

		//得到当前客户端请求的Request数据
		req := Request{
			conn: c,
			msg:  msg,
		}
		//从路由Routers中找到注册绑定 Conn的对应Handle
		go func(request siface.Request) {
			//执行注册的路由方法
			c.Router.PreHandle(request)
			c.Router.Handle(request)
			c.Router.PostHandle(request)
		}(&req)

	}
}

func (c *Connection) Start() {
	go c.StartReader() //开启处理该连接读取到客户端数据之后的请求业务

	for {
		select {
		case <-c.ExitBuffChan:
			return //得到退出消息，不再阻塞
		}
	}
}

// 停止连接，结束当前连接状态M
func (c *Connection) Stop() {
	if c.isClosed == true { //1.如果当前连接已关闭
		return
	}
	c.isClosed = true

	c.Conn.Close()         //关闭socket连接
	c.ExitBuffChan <- true //通知缓冲队列数据的业务，该连接已经关闭

	close(c.ExitBuffChan) //关闭该连接全部管道
}

func (c *Connection) GetTCPConnection() *net.TCPConn {
	return c.Conn
}

func (c *Connection) GetConnID() uint32 {
	return c.ConnID
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.Conn.RemoteAddr()
}

func (c *Connection) SendMsg(msgId uint32, data []byte) error {
	if c.isClosed == true {
		return errors.New("Connection closed when send msg")
	}

	//将data封包，并且发送
	dp := NewDataPack()
	msg, err := dp.Pack(NewMsgPackage(msgId, data))
	if err != nil {
		fmt.Println("Pack err msg id =", msgId)
		return errors.New("Pack error msg ")
	}

	//写回客户端
	if _, err := c.Conn.Write(msg); err != nil {
		fmt.Println("Write msg id", msgId, " error")
		c.ExitBuffChan <- true
		return errors.New("conn Write error")
	}
	return nil
}
