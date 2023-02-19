package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *R) ReadMessage(conn *websocket.Conn) error {
	for {
		_, message, err := conn.ReadMessage() //read received message
		if err != nil {
			return err
		}
		room.Broadcast <- message //notify the room for new message

		log.Printf("Received Message: %s\n", message)
	}
}
