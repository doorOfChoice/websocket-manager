package main

import (
	"fmt"
	"net/http"

	"../tool/wspool"

	ws "../github.com/gorilla/websocket"
)

var (
	onMessage = func(conn *wspool.Conn, message wspool.Message) {
		all := conn.GetAllConns()
		for _, v := range all {
			v.(*wspool.Conn).WriteTextMessage(message.Message)
		}
	}

	onClose = func(conn *wspool.Conn) {
		all := conn.GetAllConns()
		for _, v := range all {
			v.(*wspool.Conn).WriteTextMessage([]byte(
				fmt.Sprintf("用户 %d 离开了", conn.GetIndex()),
			))
		}
	}

	manager *wspool.Manager
)

func init() {
	manager = wspool.NewManager()
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
