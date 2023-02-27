package router

import (
	"github.com/chinathaip/chatify/chatroom"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute(h *chatroom.Hub) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.GET("/ws", handleSocket(h))

	return e
}
