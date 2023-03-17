package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/chinathaip/chatify/service"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockParticipantService struct {
	mock.Mock
}

func (m *mockParticipantService) AddAsParticipant(*service.ChatParticipants) error {
	return nil
}

func (m *mockParticipantService) Exist(chatParticipant *service.ChatParticipants) bool {
	args := m.Called(chatParticipant)
	return args.Bool(0)
}

func TestValdateJWT(t *testing.T) {
	tests := []struct {
		name        string
		token       *jwt.Token
		shouldExist bool
		statusCode  int
		message     string
	}{
		{
			name:        "Return 200 - valid token, and user has permission to view messages",
			token:       jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": uuid.NewString()}),
			shouldExist: true,
			statusCode:  http.StatusOK,
			message:     "pass the jwt auth",
		},
		{
			name:        "Return 403 - valid token, but user doesnt have permission",
			token:       jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": uuid.NewString()}),
			shouldExist: false,
			statusCode:  http.StatusForbidden,
			message:     "user has no access to view messages in this chat room",
		},
		{
			name:        "Return 401 - empty uuid",
			token:       jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": ""}),
			shouldExist: false,
			statusCode:  http.StatusUnauthorized,
			message:     "missing token",
		},
		{
			name:        "Return 401 - invalid uuid",
			token:       jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": "im not a valid uuid"}),
			shouldExist: false,
			statusCode:  http.StatusUnauthorized,
			message:     "incorrect sub value",
		},
		{
			name:        "Return 401 - no sub",
			token:       jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"Haiyah": "Fuiyoh"}),
			shouldExist: false,
			statusCode:  http.StatusUnauthorized,
			message:     "no key 'sub'",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var (
				e           = echo.New()
				req         = httptest.NewRequest(http.MethodGet, "/", nil)
				rec         = httptest.NewRecorder()
				c           = e.NewContext(req, rec)
				mockService = &mockParticipantService{}
				handler     = NewHandler(nil, nil, nil, mockService)

				tokenString, _ = test.token.SignedString([]byte(os.Getenv("JWT_SECRET")))
			)
			mockService.On("Exist", mock.AnythingOfType("*service.ChatParticipants")).Return(test.shouldExist)
			c.SetPath("/messages")
			c.SetParamNames("chat_id")
			c.SetParamValues("41")
			req.Header.Add("Authorization", "Bearer "+tokenString)

			err := handler.validateJWT(func(c echo.Context) error {
				return c.String(http.StatusOK, "pass the jwt auth")
			})(c)

			assert.NoError(t, err)
			assert.Equal(t, test.statusCode, rec.Code)
			assert.True(t, strings.ContainsAny(rec.Body.String(), test.message))
		})
	}
}
