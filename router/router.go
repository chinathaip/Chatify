package router

import (
	"github.com/chinathaip/chatify/hub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute(h *hub.H, handler *Handler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.GET("/ws", handleSocket(h))

	e.GET("/chats", handler.handleGetAllChat)
	e.POST("/chats", handler.handleCreateNewChat)
	e.GET("/messages/:chat_id", handler.handleGetMessages)
	e.POST("/messages", handler.handleStoreMessage)

	return e
}
