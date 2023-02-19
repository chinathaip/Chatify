package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *CR) ReadMessage(conn *websocket.Conn, broadcast chan []byte) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}
		broadcast <- message

		log.Printf("Received Message: %s\n", message)
	}
}
