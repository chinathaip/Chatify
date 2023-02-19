package router

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute() *echo.Echo {
	e := echo.New()
	e.Use(middleware.Logger())
	e.GET("/ws", handleSocket)

	return e
}
