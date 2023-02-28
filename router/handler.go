package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinathaip/chatify/service"
	"github.com/labstack/echo/v4"
)

type handler struct {
	chatService    service.ChatService
	messageService service.MessageService
}

func newHandler(chatService service.ChatService, messageService service.MessageService) *handler {
	return &handler{
		chatService:    chatService,
		messageService: messageService,
	}
}

func (h *handler) handleGetAllChat(c echo.Context) error {
	chat, err := h.chatService.GetAllChat()
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
		return c.String(http.StatusBadRequest, "invalid param")
	}
	msg, err := h.messageService.GetMessagesInChat(chatID)
	if len(msg) == 0 {
		log.Println("Error getting messages: ", err) //never show SQL error to client
		return c.String(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, map[string]any{"chat_id": chatID, "messages": msg})
}

func (h *handler) handleStoreMessage(c echo.Context) error {
	var msg service.Message
	if err := c.Bind(&msg); err != nil {
		log.Println("Erorr binding: ", err)
		return c.String(http.StatusBadRequest, "error!")
	}

	if msg.SenderID == 0 || msg.ChatID == 0 || msg.Data == "" {
		log.Println("Client has sent invalid request body")
		return c.String(http.StatusBadRequest, "invalid request body")
	}

	if err := h.messageService.StoreNewMessage(&msg); err != nil {
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, msg)
}
