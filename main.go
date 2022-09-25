package main

import "github.com/CharmingZhou/myserv/snet"

func main() {
	s := snet.NewServer("[myserv V0.2]")
	// 2.开启服务
	s.Serve()
}
