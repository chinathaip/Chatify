package router

import (
	"github.com/chinathaip/chatify/chatroom"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())

	room := chatroom.New()
	go room.Init()

	e.GET("/ws", handleSocket(room))

	return e
}
