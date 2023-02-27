package chatroom

import (
	"github.com/gorilla/websocket"
)

type R struct {
	name  string
	users map[*websocket.Conn]bool
}

func New(roomName string) *R {
	return &R{
		name:  roomName,
		users: make(map[*websocket.Conn]bool),
	}
}
