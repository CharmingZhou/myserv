package snet

import "github.com/CharmingZhou/myserv/siface"

type Request struct {
	conn siface.Connection //已经和客户端建立好的连接
	data []byte            //客户端请求的数据
}

func (r *Request) GetConnection() siface.Connection {
	return r.conn
}

func (r *Request) GetData() []byte {
	return r.data
}
