package ynet

import (
	"errors"
	"fmt"
	"github.com/justcy/ygo/ygo/yiface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]yiface.IConnection //管理链接信息
	connLock   sync.RWMutex                    //读写链接的读写锁
}

func (connMgr *ConnManager) Add(conn yiface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnId()] = conn

	fmt.Println("connection add to ConnManager successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn yiface.IConnection) {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	delete(connMgr.connections,conn.GetConnId())
	fmt.Println("connection Remove ConnID=",conn.GetConnId(), " successfully: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connId uint32) (yiface.IConnection, error) {
	//保护共享资源Map 加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connId]; ok {
		return conn, nil
	} else {
		return nil, errors.New("connection not found")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	//保护共享资源Map 加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	//停止并删除全部的连接信息
	for connID, conn := range connMgr.connections {
		//停止
		conn.Stop()
		//删除
		delete(connMgr.connections,connID)
	}
	fmt.Println("Clear All Connections successfully: conn num = ", connMgr.Len())
}

func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: map[uint32]yiface.IConnection{},
	}
}