package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinathaip/chatify/service"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type handler struct {
	chatservice    service.ChatService
	messageService service.MessageModel
}

func newHandler(chatService service.ChatService, db *gorm.DB) *handler {
	return &handler{
		chatservice:    chatService,
		messageService: service.MessageModel{DB: db},
	}
}

func (h *handler) handleGetAllChat(c echo.Context) error {
	chat, err := h.chatservice.GetAllChat()
	if err != nil {
		log.Println("Error retreiving chat: ", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, chat)
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
