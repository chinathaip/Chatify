package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGetAllChat(t *testing.T) {
	t.Run("Should Return 200 Ok if no error", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		c.SetPath("/chats")
		mockChatService := &mockChatService{}
		handler := NewHandler(mockChatService, nil, nil)

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
			handler := NewHandler(nil, mockMessageService, nil)

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
				"sender_id": "27326c5b-7395-435b-bfc3-330ad6686e53",
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
			handler := NewHandler(nil, mockMessageService, nil)

			handler.handleStoreMessage(c)

			assert.Equal(t, test.expectedStatusCode, c.Response().Status)
			assert.Equal(t, test.expectedServiceCalled, mockMessageService.isStoreNewMessageCalled)
		})
	}
}

func TestCreateNewUser(t *testing.T) {
	tests := []struct {
		name                 string
		userJSON             string
		expectedStatusCode   int
		expectdServiceCalled bool
	}{
		{
			name: "Should return 201 Created when valid request body",
			userJSON: `{
				"id":"27326c5b-7395-435b-bfc3-330ad6686e53",
				"username": "Kidcat"
			}`,
			expectedStatusCode:   http.StatusCreated,
			expectdServiceCalled: true,
		},
		{

			name: "Should return 404 Bad Request when invalid request body",
			userJSON: `{
				"heeeHAA: 1"
				"username": "Kidcat"
			}`,
			expectedStatusCode:   http.StatusBadRequest,
			expectdServiceCalled: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(test.userJSON))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath("/chats")
			mockUserService := &mockUserService{}
			handler := NewHandler(nil, nil, mockUserService)

			handler.handleCreateNewUser(c)

			assert.Equal(t, test.expectedStatusCode, c.Response().Status)
			assert.Equal(t, test.expectdServiceCalled, mockUserService.isCreateNewUserCalled)

		})
	}
}
