package snet

import "github.com/CharmingZhou/myserv/siface"

type BaseRouter struct{}

func (br *BaseRouter) PreHandle(req siface.Request)  {}
func (br *BaseRouter) handle(req siface.Request)     {}
func (br *BaseRouter) PostHandle(req siface.Request) {}
