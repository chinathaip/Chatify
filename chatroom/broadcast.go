package chatroom

import (
	"log"

	"github.com/gorilla/websocket"
)

func (room *R) broadcastMsg(users map[*websocket.Conn]bool, message []byte) {
	for user, active := range users {
		if active {
			err := user.WriteMessage(websocket.TextMessage, message)
			log.Println("Broadcasting to Chat Room: ", string(message))
			if err != nil {
				log.Println("Error sending to client: ", err)
				user.Close()
			}
		}
	}
}
