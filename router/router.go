package router

import (
	"net/http"

	"github.com/chinathaip/chatify/hub"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func RegRoute(h *hub.H, handler *Handler) *echo.Echo {
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	},
	))
	e.GET("/ws", handleSocket(h))

	e.GET("/chats", handler.handleGetAllChat)
	e.POST("/chats", handler.handleCreateNewChat)
	e.DELETE("/chats", handler.handleDeleteChat)
	e.GET("/chats/:chat_id/messages", handler.handleGetMessages, handler.validateJWT)

	e.POST("/users", handler.handleCreateNewUser)

	return e
}
