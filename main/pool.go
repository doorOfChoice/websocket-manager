package main

import (
	"sync"
)

//Conn连接管理池
type Pool struct {
	values map[int]interface{}
	mux    *sync.Mutex
}

func NewPool() *Pool {
	return &Pool{
		values: make(map[int]interface{}),
		mux:    &sync.Mutex{},
	}
}

//添加到连接池
//如何期间执行其他Pool操作会阻塞
func (this *Pool) Push(index int, wsConn *Conn) {
	this.mux.Lock()
	defer this.mux.Unlock()

	this.values[index] = wsConn
}

//从连接池删除
//期间执行其他Pool操作会阻塞
func (this *Pool) Delete(index int) bool {
	this.mux.Lock()
	defer this.mux.Unlock()

	if this.values[index] != nil {
		delete(this.values, index)
		return true
	}

	return false
}

//复制一个values到新的map[int]interface{}，并在期间对其他操作加锁
//防止复制期间新增和删除造成数据不同步
//期间执行其他Pool操作会阻塞
func (this *Pool) GetValues() map[int]interface{} {
	this.mux.Lock()
	defer this.mux.Unlock()

	n := make(map[int]interface{})

	for i, j := range this.values {
		n[i] = j
	}

	return n
}
