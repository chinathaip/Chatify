package router

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/chinathaip/chatify/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func (h *Handler) validateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Missing token")
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid Tokenn")
		}

		cid := c.Param("chat_id")
		chatID, err := strconv.Atoi(cid)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid Param")
		}

		userID := claims["sub"].(string)
		uid, err := uuid.Parse(userID)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		if ok := h.participantService.Exist(&service.ChatParticipants{ChatID: chatID, UserID: uid}); !ok {
			return c.JSON(http.StatusForbidden, "user has no access to view messages in this chat room")
		}

		c.Set("claims", claims)

		return next(c)
	}
}
