package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/chinathaip/chatify/service"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

type mockChatService struct {
	isGetAllCalled bool
}

func (cs *mockChatService) GetAllChat() ([]service.Chat, error) {
	cs.isGetAllCalled = true
	return []service.Chat{
		{
			ID:   1,
			Name: "Test Chat Room",
		},
		{
			ID:   2,
			Name: "Test Chat Room 2",
		},
	}, nil
}

func TestGetAllChat(t *testing.T) {
	t.Run("Should Return 200 Ok if no error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/chats")
		handler := newHandler(&mockChatService{}, nil)
		expected := []service.Chat{
			{
				ID:   1,
				Name: "Test Chat Room",
			},
			{
				ID:   2,
				Name: "Test Chat Room 2",
			},
		}

		handler.handleGetAllChat(c)

		result, _ := handler.chatservice.GetAllChat()
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, expected, result)
	})
}
