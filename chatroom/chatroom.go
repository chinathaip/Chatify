package chatroom

import (
	"github.com/gorilla/websocket"
)

type R struct {
	name       string
	users      map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func New(roomName string) *R {
	return &R{
		name:       roomName,
		users:      make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
	}
}

func (room *R) Init() {
	for {
		select {
		case message := <-room.Broadcast: //broadcast new message when notified
			room.broadcastMsg(room.users, message)
		case client := <-room.Register: //register new user when notified
			room.users[client] = true
		case client := <-room.Unregister: //unregister user when notified
			room.users[client] = false
		}
	}
}
