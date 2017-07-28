package main

import (
	"encoding/json"
)

//通讯信息类
//这个类负责在客户端和服务器之间传递消息
type Message struct {
	To          int
	From        int
	Message     string
	Command     string
	MessageType int
}

//转化成JSON字符串
func (this *Message) ToJson() ([]byte, error) {
	return json.Marshal(*this)
}

//转化为对象
func ParseMessage(buf []byte) Message {
	msg := Message{}
	json.Unmarshal(buf, &msg)

	return msg
}
