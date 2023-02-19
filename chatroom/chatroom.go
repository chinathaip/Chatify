package chatroom

import (
	"github.com/gorilla/websocket"
)

type CR struct {
	name       string
	users      map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func New(roomName string) *CR {
	return &CR{
		name:       roomName,
		users:      make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (room *CR) Init() {
	for {
		select {
		case message := <-room.Broadcast:
			room.broadcastMsg(room.users, message)
		case client := <-room.Register:
			room.users[client] = true
		case client := <-room.Unregister:
			room.users[client] = false
		}
	}
}
