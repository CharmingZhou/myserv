package siface

// 路由接口，这里路由是使用框架者给该连接自定的处理业务方法，路由里的Request则包含用该链接的连接信息和该链接的请求数据信息
type Router interface {
	PreHandle(request Request)  //在处理conn业务前的钩子方法
	Handle(request Request)     //处理conn业务的方法
	PostHandle(request Request) //处理conn业务之后的钩子方法
}
