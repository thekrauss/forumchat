package wsk

import (
	"github.com/gorilla/websocket"
)

type UserChat struct {
	channels   *Channel
	Username   string
	Connection *websocket.Conn
}

func NewUserChat(channels *Channel, username string, conn *websocket.Conn) *UserChat {
	return &UserChat{
		channels:   channels,
		Username:   username,
		Connection: conn,
	}
}
