package chatroom

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Run("Server receives and broadcasts message correctly", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		room := New("Test Room", ctx)
		go room.Init()
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			upgrader := websocket.Upgrader{}
			conn, err := upgrader.Upgrade(w, r, nil)
			assert.NoError(t, err)
			room.Register <- conn
			go room.ReadMsgFrom(conn)
		}))
		defer server.Close()

		wsURL := "ws" + server.URL[4:]

		user1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
		assert.NoError(t, err)

		expected := []byte("Test Message")
		err = user1.WriteMessage(websocket.TextMessage, expected)
		assert.NoError(t, err)

		assert.Equal(t, expected, <-room.Broadcast)
	})

}
