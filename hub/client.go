package hub

import (
	"net"
)

type Client struct {
	roomName string
	conn     connection
}

type connection interface {
	RemoteAddr() net.Addr
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, data []byte, err error)
}

func NewClient(roomName string, conn connection) *Client {
	return &Client{
		roomName: roomName,
		conn:     conn,
	}
}

type Message struct {
	roomName string
	data     *JSONMessage
}
