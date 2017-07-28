package main

import (
	"fmt"
	"net/http"

	ws "../github.com/gorilla/websocket"
)

var (
	onMessage = func(conn *Conn, message Message) {
		all := conn.GetAllConns()
		for _, v := range all {
			v.(*Conn).WriteTextMessage([]byte(message.Message))
		}
	}

	onClose = func(conn *Conn) {
		all := conn.GetAllConns()
		for _, v := range all {
			v.(*Conn).WriteTextMessage([]byte(
				fmt.Sprintf("用户 %d 离开了", conn.GetIndex()),
			))
		}
	}
	manager *Manager
)

func init() {
	manager = NewManager()
	manager.OnClose = onClose
	manager.OnMessage = onMessage
}

func handle(w http.ResponseWriter, r *http.Request) {
	up := &ws.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	wsConn, err := up.Upgrade(w, r, nil)

	if err != nil {
		return
	}

	manager.Push(wsConn)
}

func main() {
	http.HandleFunc("/ws", handle)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "../test.html")
	})

	http.ListenAndServe(":8888", nil)
}
