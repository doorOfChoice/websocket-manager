package main

import (
	"sync"

	ws "../github.com/gorilla/websocket"
)

//自增计数，为链接自动命名
var auto = 0

//管理器
type Manager struct {
	mux       *sync.Mutex          //锁
	pool      *Pool                //链接管理池
	OnMessage func(*Conn, Message) //客户端收到信息的函数模板
	OnClose   func(*Conn)          //客户端关闭链接的函数模板
}

func NewManager() *Manager {
	return &Manager{
		mux:  &sync.Mutex{},
		pool: NewPool(),
	}
}

//添加新的WebSocket到连接池
//WebSocket链接会转化为Conn链接
func (this *Manager) Push(conn *ws.Conn) *Conn {
	this.mux.Lock()
	defer this.mux.Unlock()

	c := NewConn(auto, conn, this)

	c.bindOnMessage(this.OnMessage)
	c.bindOnClose(this.OnClose)

	this.pool.Push(auto, c)

	//开始监听Conn
	go c.startListen()
	//计数自增
	auto++

	return c
}

//从连接池删除指定连接
func (this *Manager) Delete(index int) bool {
	return this.pool.Delete(index)
}

//获取全部的在线的连接
func (this *Manager) GetAllConns() map[int]interface{} {
	return this.pool.GetValues()
}
