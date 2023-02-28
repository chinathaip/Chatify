package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

func (cs *mockChatService) DeleteChat(chatID int) error {
	return nil
}

func (cs *mockChatService) CreateNewChat(*service.Chat) error {
	return nil
}

func (cs *mockChatService) IsChatExist(chatName string) (int, bool) {
	return 0, true
}

type mockMessageService struct {
	isGetMessagesCalled     bool
	isStoreNewMessageCalled bool
}

type mockError struct{}

func (e *mockError) Error() string {
	return "Error occured!"
}

func (ms *mockMessageService) GetMessagesInChat(chatID int) ([]service.Message, error) {
	ms.isGetMessagesCalled = true
	if chatID == 1 {
		return []service.Message{
			{
				ID:       1,
				SenderID: 1,
				ChatID:   1,
				Data:     "Message 1",
			},
			{
				ID:       2,
				SenderID: 1,
				ChatID:   1,
				Data:     "Message 2 from the same dude",
			},
		}, nil
	}
	return []service.Message{}, &mockError{}
}

func (ms *mockMessageService) StoreNewMessage(msg *service.Message) error {
	ms.isStoreNewMessageCalled = true
	return nil
}

func TestGetAllChat(t *testing.T) {
	t.Run("Should Return 200 Ok if no error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/chats")
		mockChatService := &mockChatService{}
		handler := NewHandler(mockChatService, nil)

		err := handler.handleGetAllChat(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.True(t, mockChatService.isGetAllCalled)
	})
}

func TestGetMessages(t *testing.T) {
	tests := []struct {
		name                  string
		paramValue            string
		expectedStatusCode    int
		expectedServiceCalled bool
	}{
		{
			name:                  "Should return 200 Ok if chat_id is valid",
			paramValue:            "1",
			expectedStatusCode:    http.StatusOK,
			expectedServiceCalled: true,
		},
		{
			name:                  "Should return 400 Bad Request if chat_id is not valid (enter 0)",
			paramValue:            "0",
			expectedStatusCode:    http.StatusBadRequest,
			expectedServiceCalled: false,
		},
		{
			name:                  "Should return 400 Bad Request if chat_id is not valud (enter string)",
			paramValue:            "dadoadjiw",
			expectedStatusCode:    http.StatusBadRequest,
			expectedServiceCalled: false,
		},
		{
			name:                  "Should return 404 Not Found if chat_id doesn't exist in db",
			paramValue:            "999",
			expectedStatusCode:    http.StatusNotFound,
			expectedServiceCalled: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/messages")
			c.SetParamNames("chat_id")
			c.SetParamValues(test.paramValue)
			mockMessageService := &mockMessageService{}
			handler := NewHandler(nil, mockMessageService)

			handler.handleGetMessages(c)

			assert.Equal(t, test.expectedStatusCode, c.Response().Status)
			assert.Equal(t, test.expectedServiceCalled, mockMessageService.isGetMessagesCalled)

		})
	}
}

func TestStoreMessage(t *testing.T) {

	tests := []struct {
		name                  string
		msgJSON               string
		expectedStatusCode    int
		expectedServiceCalled bool
	}{
		{
			name: "Should return 201 Created when valid request body",
			msgJSON: `{
				"sender_id": 1,
				"chat_id": 1,
				"data": "This message was created by API"
			}`,
			expectedStatusCode:    http.StatusCreated,
			expectedServiceCalled: true,
		},
		{
			name: "Should return 400 Bad Request when invalid request body",
			msgJSON: `{
				"ayoyoyoyo" : 1
			}`,
			expectedStatusCode:    http.StatusBadRequest,
			expectedServiceCalled: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.msgJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/messages")
			mockMessageService := &mockMessageService{}
			handler := NewHandler(nil, mockMessageService)

			handler.handleStoreMessage(c)

			assert.Equal(t, test.expectedStatusCode, c.Response().Status)
			assert.Equal(t, test.expectedServiceCalled, mockMessageService.isStoreNewMessageCalled)
		})
	}
}
