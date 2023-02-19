package socket

import (
	"log"

	"github.com/gorilla/websocket"
)

func ReadLoops(conn *websocket.Conn) error {
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			return err
		}

		log.Printf("Received Message: %s\n", message)
	}
}
