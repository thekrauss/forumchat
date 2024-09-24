package wsk

import "time"

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

type Message struct {
	ID             int       `json:"id"`
	SenderUsername string    `json:"senderUsername"`
	TargetUsername string    `json:"targetUsername"`
	Content        string    `json:"content"`
	Timestamp      time.Time `json:"timestamp"`
	IsTyping       bool      `json:"isTyping"`
	Type           string    `json:"type"`
}

type userChannel chan *UserChat
type messageChannel chan *Message

type Channel struct {
	messageChannel messageChannel
	leaveChannel   userChannel
}
