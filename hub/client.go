package hub

import (
	"net"

	"github.com/google/uuid"
)

type Client struct {
	roomName string
	userID   uuid.UUID
	conn     connection
}

type connection interface {
	RemoteAddr() net.Addr
	WriteMessage(messageType int, data []byte) error
	ReadMessage() (messageType int, data []byte, err error)
}

func NewClient(roomName string, userID uuid.UUID, conn connection) *Client {
	return &Client{
		roomName: roomName,
		userID:   userID,
		conn:     conn,
	}
}

type Message struct {
	roomName string
	data     *JSONMessage
}
