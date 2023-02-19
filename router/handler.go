package router

import (
	"net/http"
	"sync"

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

func handleSocket(rooms *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomName := c.QueryParam("roomName")
		if roomName == "" {
			roomName = "Test Chat Room"
		}

		data, _ := rooms.LoadOrStore(roomName, chatroom.New(roomName))

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		room := data.(*chatroom.R)
		go room.Init()

		room.Register <- conn //notify room when user join
		room.ReadMessage(conn)
		room.Unregister <- conn //notify room when user left
		return nil
	}
}
