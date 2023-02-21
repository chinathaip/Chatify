package chatroom

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type R struct {
	name       string
	users      map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
	Ctx        context.Context
}

func New(roomName string, ctx context.Context) *R {
	return &R{
		name:       roomName,
		users:      make(map[*websocket.Conn]bool),
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Ctx:        ctx,
	}
}

func (room *R) Init() {
	for {
		select {
		case message := <-room.Broadcast: //broadcast new message when notified
			room.broadcastMsg(room.users, message)
		case client := <-room.Register: //register new user when notified
			room.users[client] = true
			log.Println("New client: ", client.RemoteAddr())
		case client := <-room.Unregister: //unregister user when notified
			delete(room.users, client)
		case <-room.Ctx.Done():
			log.Println("Terminating Chat rooom: ", room.name)
			return
		}
	}
}

func (room *R) MonitorUser(cancel context.CancelFunc) {
	for {
		if len(room.users) <= 0 {
			cancel()
		}
	}
}
