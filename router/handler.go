package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinathaip/chatify/service"
	"github.com/labstack/echo/v4"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type handler struct {
	messageService service.MessageModel
}

func newHandler(dsn string) *handler {
	gorm, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalln("Error connecting to db: ", err)
	}
	return &handler{
		messageService: service.MessageModel{DB: gorm},
	}
}
func (h *handler) handleGetMessages(c echo.Context) error {
	id := c.Param("chat_id")
	chatID, err := strconv.Atoi(id)
	if err != nil || chatID == 0 {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	msg, err := h.messageService.GetMessagesInChat(chatID)
	if len(msg) == 0 {
		log.Println("Error getting messages: ", err) //never show SQL error to client
		return echo.ErrNotFound
	}
	return c.JSON(http.StatusOK, map[string]any{"chat_id": chatID, "messages": msg})
}

func (h *handler) handleStoreMessage(c echo.Context) error {
	var msg service.Message
	if err := c.Bind(&msg); err != nil {
		log.Println("Erorr binding: ", err)
		return echo.ErrBadRequest
	}

	if err := h.messageService.StoreNewMessage(&msg); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, msg)
}
