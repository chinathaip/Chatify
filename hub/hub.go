package hub

import (
	"context"
	"log"
	"sync"

	"github.com/chinathaip/chatify/service"
	"github.com/gorilla/websocket"
)

type H struct {
	Broadcast   chan *Message
	Register    chan *Client
	Unregister  chan *Client
	Rooms       map[string]*Room
	mutex       sync.RWMutex
	msgService  service.MessageService
	chatService service.ChatService
}

func New(chatService service.ChatService, msgService service.MessageService) *H {
	return &H{
		Broadcast:   make(chan *Message),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Rooms:       make(map[string]*Room),
		mutex:       sync.RWMutex{},
		msgService:  msgService,
		chatService: chatService,
	}
}

func (h *H) setNewRoom(roomName string, r *Room) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.Rooms[roomName] = r
}

func (h *H) getRoom(roomName string) *Room {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.Rooms[roomName]
}

func (h *H) deleteRoom(roomName string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	delete(h.Rooms, roomName)
}

func (h *H) Init(ctx context.Context) {
run:
	for {
		select {
		case client := <-h.Register:
			room := h.getRoom(client.roomName)
			if room == nil { //if room doesnt exist in hub
				room = NewRoom(client.roomName)
			}
			room.setNewUser(client.conn)
			h.setNewRoom(client.roomName, room)

			//store in db
			if h.chatService == nil {
				continue
			}
			if id, exist := h.chatService.IsChatExist(client.roomName); exist {
				room.id = id
				continue
			}
			chat := &service.Chat{Name: client.roomName}
			err := h.chatService.CreateNewChat(chat)
			if err != nil {
				log.Panicln("Error creating new chat: ", err)
			}
			room.id = chat.ID

		case client := <-h.Unregister:
			room := h.getRoom(client.roomName)
			if room != nil {
				room.deleteUser(client.conn)
				if len(room.users) == 0 { //if no one is left in the room
					h.deleteRoom(client.roomName)
				}
			}

		case message := <-h.Broadcast:
			room := h.getRoom(message.roomName) //get room to send the message
			if room != nil {
				for user, active := range room.users {
					if !active {
						return
					}
					err := user.WriteMessage(websocket.TextMessage, message.data) //send the message
					if err != nil {
						log.Panicf("Error broadcasting message from %s", user.RemoteAddr())
					}
					log.Printf("Broadcasting to : %s with message %s", user.RemoteAddr(), message.data)

				}

				//store in db
				if h.msgService == nil {
					continue
				}
				msg := &service.Message{SenderID: 1, ChatID: room.id, Data: string(message.data)}
				err := h.msgService.StoreNewMessage(msg)
				if err != nil {
					log.Panicln("Error storing messages: ", err)
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
				break
			}
		}
		message := &Message{roomName: client.roomName, data: data}
		log.Printf("Client: %s sent : %s\n", client.conn.RemoteAddr(), string(message.data))

		h.Broadcast <- message //send for broadcasting
	}
}
