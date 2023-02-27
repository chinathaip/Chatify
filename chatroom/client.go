package chatroom

import "github.com/gorilla/websocket"

type Client struct {
	roomName string
	conn     *websocket.Conn
}

func NewClient(roomName string, conn *websocket.Conn) *Client {
	return &Client{
		roomName: roomName,
		conn:     conn,
	}
}

type Message struct {
	roomName string
	data     []byte
}
