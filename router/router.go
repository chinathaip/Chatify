package router

import (
	"github.com/chinathaip/chatify/hub"
	"github.com/chinathaip/chatify/service"
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
