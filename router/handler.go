package router

import (
	"net/http"

	"github.com/chinathaip/chatify/chatroom"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func handleSocket(h *chatroom.Hub) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomName := c.QueryParam("roomName")
		if roomName == "" {
			roomName = "Test Chat Room"
		}
		// data, _ := rooms.LoadOrStore(roomName, chatroom.New(roomName, ctx))

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := chatroom.NewClient(roomName, conn)

		h.Register <- client
		h.ReadMsgFrom(client) //read message from client
		return nil
	}
}
