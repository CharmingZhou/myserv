package siface

type Server interface {
	Start()
	Stop()
	Serve()
}
