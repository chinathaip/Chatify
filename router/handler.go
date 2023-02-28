package router

import (
	"log"
	"net/http"
	"strconv"

	"github.com/chinathaip/chatify/service"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	chatService    service.ChatService
	messageService service.MessageService
}

func NewHandler(chatService service.ChatService, messageService service.MessageService) *Handler {
	return &Handler{
		chatService:    chatService,
		messageService: messageService,
	}
}

func (h *Handler) handleGetAllChat(c echo.Context) error {
	chat, err := h.chatService.GetAllChat()
	if err != nil {
		log.Println("Error retreiving chat: ", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, chat)
}

func (h *Handler) handleCreateNewChat(c echo.Context) error {
	var chat service.Chat
	err := c.Bind(&chat)
	if err != nil || chat.Name == "" {
		log.Println("Client has sent invalid request body")
		return c.String(http.StatusBadRequest, "invalid param")
	}

	err = h.chatService.CreateNewChat(&chat)
	if err != nil {
		log.Println("Internal error :", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, chat)
}

func (h *Handler) handleGetMessages(c echo.Context) error {
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

func (h *Handler) handleStoreMessage(c echo.Context) error {
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
		log.Println("Internal error :", err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, msg)
}
