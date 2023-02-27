package hub

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
)

type H struct {
	Broadcast  chan *Message
	Register   chan *Client
	Unregister chan *Client
	Rooms      map[string]*Room
}

func New() *H {
	return &H{
		Broadcast:  make(chan *Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Rooms:      make(map[string]*Room),
	}
}

func (h *H) Init(ctx context.Context) {
run:
	for {
		select {
		case client := <-h.Register:
			room := h.Rooms[client.roomName] //get room from hub
			if room == nil {                 //if room doesnt exist
				room = NewRoom(client.roomName) //create new room
			}
			room.users[client.conn] = true  //add client to the room
			h.Rooms[client.roomName] = room //add room to the hub

		case client := <-h.Unregister:
			room := h.Rooms[client.roomName]
			if room != nil {
				delete(room.users, client.conn) //delete client from room
				if len(room.users) == 0 {       //if no one is left in the room
					delete(h.Rooms, client.roomName) //delete the room from the hub
				}
			}
		case message := <-h.Broadcast:
			room := h.Rooms[message.roomName] //get room to send the message
			if room != nil {
				for user, active := range room.users {
					if !active {
						return
					}
					err := user.WriteMessage(websocket.TextMessage, message.data) //send the message
					if err != nil {
						log.Printf("Error broadcasting message from %s", user.RemoteAddr())
					}
					log.Printf("Broadcasting to : %s with message %s", user.RemoteAddr(), message.data)
				}
			}
		case <-ctx.Done():
			break run
		}
	}
}

func (h *H) ReadMsgFrom(client *Client) {
	for {
		_, data, err := client.conn.ReadMessage() //read received message
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Client %s has closed connection", client.conn.RemoteAddr())
				h.Unregister <- client
				return
			}
		}
		message := &Message{roomName: client.roomName, data: data}
		log.Printf("Client: %s sent : %s\n", client.conn.RemoteAddr(), string(message.data))

		h.Broadcast <- message //send for broadcasting
	}
}
