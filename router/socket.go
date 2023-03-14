package router

import (
	"errors"
	"net/http"

	"github.com/chinathaip/chatify/hub"
	"github.com/google/uuid"
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

func handleSocket(h *hub.H) echo.HandlerFunc {
	return func(c echo.Context) error {
		roomName := c.QueryParam("roomName")
		if roomName == "" {
			roomName = "Test Chat Room"
		}

		userID := c.QueryParam("userID")
		uuid, err := uuid.Parse(userID)
		if err != nil {
			return errors.New("invalid UUID")
		}

		//upgrade connection to web socket
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := hub.NewClient(roomName, uuid, conn)

		h.Register <- client
		h.ReadMsgFrom(client) //read message from client
		return nil
	}
}
