package snet

import (
	"fmt"
	"net"

	"github.com/CharmingZhou/myserv/siface"
)

type Connection struct {
	Conn     *net.TCPConn //当前连接的socket TCP套接字
	ConnID   uint32       //当前连接的ID 也可以作为SessionID，ID全局唯一
	isClosed bool         //当前连接的关闭状态

	handleAPI siface.HandFunc //该连接的处理方法api

	ExitBuffChan chan bool //告知该连接已经退出/停止的channel
}

// 创建连接的方法
func NewConnection(conn *net.TCPConn, connID uint32, callbackApi siface.HandFunc) *Connection {
	c := &Connection{
		Conn:         conn,
		ConnID:       connID,
		isClosed:     false,
		handleAPI:    callbackApi,
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
		//读取我们最大的数据到buf中
		buf := make([]byte, 512)
		cnt, err := c.Conn.Read(buf)
		if err != nil {
			fmt.Println("recv buf err", err)
			c.ExitBuffChan <- true
			continue
		}
		//调用当前连接业务（这里执行的是当前conn的绑定的handle方法）
		if err := c.handleAPI(c.Conn, buf, cnt); err != nil {
			fmt.Println("connId ", c.ConnID, "handle is error")
			c.ExitBuffChan <- true
			return
		}
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
