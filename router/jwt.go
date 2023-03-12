package router

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
)

type jwtClaims struct {
	ChatID int
	jwt.RegisteredClaims
}

func validateJWT(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		idstr := c.Param("chat_id")
		id, err := strconv.Atoi(idstr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, "Invalid Param")
		}

		tokenString := c.Request().Header.Get("Authorization")
		if tokenString == "" {
			return c.String(http.StatusUnauthorized, "Missing token")
		}

		tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
		token, err := jwt.ParseWithClaims(tokenString, &jwtClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method")
			}
			return []byte("secret"), nil
		})
		if err != nil {
			return c.JSON(http.StatusUnauthorized, "Invalid Token")
		}

		claims, ok := token.Claims.(*jwtClaims)
		if !ok || !token.Valid {
			return c.JSON(http.StatusUnauthorized, "Invalid Tokenn")
		}

		if claims.ChatID != id {
			return c.JSON(http.StatusUnauthorized, "Invalid Tokennn")
		}

		c.Set("claims", claims)

		return next(c)
	}
}
