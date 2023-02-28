package router

import (
	"net/http"

	"github.com/chinathaip/chatify/hub"
	"github.com/chinathaip/chatify/service"
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"gorm.io/gorm"
)

func RegRoute(h *hub.H, db *gorm.DB) *echo.Echo {
	e := echo.New()
	handler := newHandler(&service.ChatModel{DB: db}, &service.MessageModel{DB: db})

	e.Use(middleware.Logger())
	e.GET("/ws", handleSocket(h))

	e.GET("/chats", handler.handleGetAllChat)
	e.GET("/messages/:chat_id", handler.handleGetMessages)
	e.POST("/messages", handler.handleStoreMessage)

	return e
}

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

		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer conn.Close()

		client := hub.NewClient(roomName, conn)

		h.Register <- client
		h.ReadMsgFrom(client) //read message from client
		return nil
	}
}
