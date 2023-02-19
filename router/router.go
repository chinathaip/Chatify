package router

import (
	"sync"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute() *echo.Echo {
	e := echo.New()
	rooms := new(sync.Map)
	e.Use(middleware.Logger())
	e.GET("/ws", handleSocket(rooms))

	return e
}
