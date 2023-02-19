package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *CR) broadcastMsg(users []*websocket.Conn, message []byte) {
	for _, user := range users {
		err := user.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			log.Println("Error sending to client: ", err)
			user.Close()
		}
	}
}
