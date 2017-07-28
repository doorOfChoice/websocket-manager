package wspool

import (
	"log"

	ws "../../github.com/gorilla/websocket"
)

//封装好的包含websocket链接的对象
type Conn struct {
	index     int
	wsConn    *ws.Conn
	manager   *Manager
	onMessage func(*Conn, Message)
	onClose   func(*Conn)
}

//新建一个Conn
func NewConn(index int, wsConn *ws.Conn, manager *Manager) *Conn {
	return &Conn{
		index:   index,
		wsConn:  wsConn,
		manager: manager,
	}
}

//获取当前链接的序号
func (this *Conn) GetIndex() int {
	return this.index
}

//获取所有在线的链接
func (this *Conn) GetAllConns() map[int]interface{} {
	return this.manager.GetAllConns()
}

//写入文本信息
func (this *Conn) WriteTextMessage(msg []byte) {
	//this.wsConn.SetWriteDeadline(time.Now().Add(5e9))
	// msgObj := Message{
	// 	To:          -1,
	// 	From:        this.GetIndex(),
	// 	Message:     string(msg),
	// 	MessageType: ws.TextMessage,
	// }
	// msgObjStr, err := msgObj.ToJson()

	// if err != nil {
	// 	log.Printf("WriteTextMessage error is %s", err.Error())
	// }

	this.wsConn.WriteMessage(ws.TextMessage, msg)
}

//读入数据
func (this *Conn) ReadMessage() (int, []byte, error) {
	//this.wsConn.SetReadDeadline(time.Now().Add(5e9))
	return this.wsConn.ReadMessage()
}

//绑定一个接受到信息后处理的函数
func (this *Conn) bindOnMessage(callback func(*Conn, Message)) {
	this.onMessage = callback
}

//绑定一个关闭前处理的函数
func (this *Conn) bindOnClose(callback func(conn *Conn)) {
	this.onClose = callback
}

//启动goroutine
//开始监听指定Conn
func (this *Conn) startListen() {
	defer func() {
		if err := recover(); err != nil {
			this.wsConn.Close()
			this.manager.Delete(this.index)
			this.onClose(this)
		}
	}()

	log.Printf("start listening %d\n", this.GetIndex())

	for {

		t, msg, err := this.wsConn.ReadMessage()

		if err != nil {
			log.Printf("startListen error is %s", err.Error())
			continue
		}

		msgObj := Message{
			Message:     msg,
			MessageType: t,
		}
		this.onMessage(this, msgObj)

	}
}
