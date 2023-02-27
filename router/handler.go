package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinathaip/chatify/db"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	messageModel db.MessageModel
}

func NewHandler(dsn string) *Handler {
	gorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error connecting to db: ", err)
	}
	return &Handler{
		messageModel: db.MessageModel{DB: gorm},
	}
}
func (h *Handler) handleGetMessages(c echo.Context) error {
	id := c.Param("chat_id")
	chatID, err := strconv.Atoi(id)
	if err != nil || chatID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	msg, err := h.messageModel.GetMessagesInChat(chatID)
	if len(msg) == 0 {
		log.Println("Error getting messages: ", err) //never show SQL error to client
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, map[string]any{"chat_id": 1, "messages": msg})
}
