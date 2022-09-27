package snet

import (
	"errors"
	"fmt"
	"sync"

	"github.com/CharmingZhou/myserv/siface"
)

type ConnManager struct {
	connections map[uint32]siface.Connection //管理连接信息
	connLock    sync.RWMutex                 //读写连接的读写锁
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]siface.Connection),
	}
}

func (connMgr *ConnManager) Add(conn siface.Connection) {
	//保护共享资源Map加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnID()] = conn
	fmt.Println("connection add to ConnManager successfully: conn num =", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn siface.Connection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//删除连接信息
	delete(connMgr.connections, conn.GetConnID())
	fmt.Println("connection Remove ConnID=", conn.GetConnID(), " successfully:conn num=", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (siface.Connection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

// 清除并停止所有连接
func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All connectins successfully: conn num =", connMgr.Len())
}
