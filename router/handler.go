package router

import (
	"log"
	"net/http"
	"strconv"
	"time"

	er "github.com/chinathaip/chatify/error"
	"github.com/chinathaip/chatify/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

var handErr = er.HandlerError{}

type Handler struct {
	chatService    service.ChatService
	messageService service.MessageService
	userService    service.UserService
}

func NewHandler(chatService service.ChatService, messageService service.MessageService, userService service.UserService) *Handler {
	return &Handler{
		chatService:    chatService,
		messageService: messageService,
		userService:    userService,
	}
}

func (h *Handler) handleGetAllChat(c echo.Context) error {
	chat, err := h.chatService.GetAllChat()
	if err != nil {
		handErr.Log(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusOK, chat)
}

func (h *Handler) handleCreateNewChat(c echo.Context) error {
	var chat service.Chat
	err := c.Bind(&chat)
	if err != nil || chat.Name == "" {
		handErr.Log(err)
		return c.String(http.StatusBadRequest, "invalid param")
	}

	err = h.chatService.CreateNewChat(&chat)
	if err != nil {
		handErr.Log(err)
		return echo.ErrInternalServerError
	}

	claims := &jwtClaims{
		chat.ID,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 1)),
		}}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if err != nil {
		handErr.Log(err)
		return echo.ErrInternalServerError

	}
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		handErr.Log(err)
		return echo.ErrBadGateway
	}

	log.Println("Token: ", t)

	return c.JSON(http.StatusCreated, chat)
}

func (h *Handler) handleDeleteChat(c echo.Context) error {
	param := c.QueryParam("chat_id")
	id, err := strconv.Atoi(param)
	if err != nil || id == 0 {
		handErr.Log(err)
		log.Println("Client has sent invalid request body")
		return c.String(http.StatusBadRequest, "invalid param")
	}

	err = h.chatService.DeleteChat(id)
	if err != nil {
		handErr.Log(err)
		return c.String(http.StatusInternalServerError, "something wrong on our end")
	}

	return c.JSON(http.StatusOK, map[string]string{"status": "deleted"})
}

func (h *Handler) handleGetMessages(c echo.Context) error {
	id := c.Param("chat_id")
	chatID, err := strconv.Atoi(id)
	if err != nil || chatID == 0 {
		return c.String(http.StatusBadRequest, "invalid param")
	}
	msg, err := h.messageService.GetMessagesInChat(chatID)
	if err != nil {
		handErr.Log(err)
	}

	if len(msg) == 0 {
		return c.String(http.StatusNotFound, "not found")
	}
	return c.JSON(http.StatusOK, map[string]any{"chat_id": chatID, "messages": msg})
}

func (h *Handler) handleStoreMessage(c echo.Context) error {
	var msg service.Message
	if err := c.Bind(&msg); err != nil {
		handErr.Log(err)
		return c.String(http.StatusBadRequest, "error!")
	}

	if msg.ChatID == 0 || msg.Data == "" {
		log.Println("Client has sent invalid request body")
		return c.String(http.StatusBadRequest, "invalid request body")
	}

	if err := h.messageService.StoreNewMessage(&msg); err != nil {
		handErr.Log(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, msg)
}

func (h *Handler) handleCreateNewUser(c echo.Context) error {
	var user service.User
	if err := c.Bind(&user); err != nil {
		handErr.Log(err)
		return c.String(http.StatusBadRequest, "error!")
	}

	if user.ID == uuid.Nil || user.Username == "" {
		log.Println("Client has sent invalid request body")
		return c.String(http.StatusBadRequest, "invalid request body")
	}

	if err := h.userService.CreateNewUser(user); err != nil {
		handErr.Log(err)
		return echo.ErrInternalServerError
	}

	return c.JSON(http.StatusCreated, user)
}
