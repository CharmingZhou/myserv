package utils

import (
	"encoding/json"
	"io/ioutil"

	"github.com/CharmingZhou/myserv/siface"
)

type GlobalObj struct {
	TcpServer siface.Server
	Host      string
	TcpPort   int
	Name      string
	Version   string

	MaxPacketSize uint32
	MaxConn       int
}

var GlobalObject *GlobalObj

func (g *GlobalObj) Reload() {
	data, err := ioutil.ReadFile("conf/myserv.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

func init() {
	GlobalObject = &GlobalObj{
		Name:          "MyServApp",
		Version:       "V0.4",
		TcpPort:       7777,
		Host:          "0.0.0.0",
		MaxConn:       12000,
		MaxPacketSize: 4096,
	}
	GlobalObject.Reload()
}
