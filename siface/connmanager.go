package siface

// 连接管理抽象层
type ConnManager interface {
	Add(conn Connection)                   //添加连接
	Remove(conn Connection)                //删除连接
	Get(connID uint32) (Connection, error) //利用ConnID 获取连接
	Len() int                              //获取当前连接
	ClearConn()                            //删除并停止所有连接
}
