package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *R) ReadMsgFrom(conn *websocket.Conn) error {
	for {
		_, message, err := conn.ReadMessage() //read received message
		if err != nil {
			return err
		}
		log.Printf("Client: %s sent : %s\n", conn.RemoteAddr().String(), message)

		room.Broadcast <- message
		log.Printf("Send for broadcasting: %s\n", message)
	}

}
