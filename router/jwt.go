package router

import (
	"fmt"
	"log"
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
		cid := c.Param("chat_id")
		chatID, err := strconv.Atoi(cid)
		if err != nil {
			log.Printf("Cannot convert param chat_id to integer: %v", err)
			return c.JSON(http.StatusBadRequest, "Invalid Param")
		}

		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			log.Println("No Token was provided")
			return c.String(http.StatusUnauthorized, "missing token")
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil {
			log.Printf("Error while Parsing jwt: %v", err)
			return c.JSON(http.StatusUnauthorized, "invalid token")
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			log.Printf("cannot validate the token - ok? : %v, valid?: %v", ok, token.Valid)
			return c.JSON(http.StatusUnauthorized, "cannot validate the token")
		}

		if _, found := claims["sub"]; !found {
			return c.JSON(http.StatusUnauthorized, "no key 'sub'")
		}

		userID := claims["sub"].(string)
		uid, err := uuid.Parse(userID)
		if err != nil {
			log.Printf("Error while Parsing string to UUID: %v", err)
			return c.JSON(http.StatusUnauthorized, "incorrect sub value")
		}

		if ok := h.participantService.Exist(&service.ChatParticipants{ChatID: chatID, UserID: uid}); !ok {
			return c.JSON(http.StatusForbidden, "user has no access to view messages in this chat room")
		}

		c.Set("claims", claims)

		return next(c)
	}
}
