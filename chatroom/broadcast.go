package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *CR) broadcastMsg(users map[*websocket.Conn]bool, message []byte) {
	for user, active := range users {
		if active {
			err := user.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println("Error sending to client: ", err)
				user.Close()
			}
		}
	}
}
