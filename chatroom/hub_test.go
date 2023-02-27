package chatroom

import (
	"context"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {

	t.Run("Register new client should create new room if not exists", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		client := NewClient("Test room", &websocket.Conn{})

		h.Register <- client
		_, found := h.Rooms["Test room"]

		assert.Equal(t, 1, len(h.Rooms))
		assert.True(t, found)
	})

	t.Run("Register new client should add new user to the existing room", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		client2 := NewClient("Room1", &websocket.Conn{})
		h.Rooms["Room1"] = New("Room1")

		h.Register <- client2
		_, found1 := h.Rooms["Room1"]

		assert.Equal(t, 1, len(h.Rooms))
		assert.True(t, found1)
	})

	t.Run("Unregister should remove user from the existing room", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		client1 := NewClient("Room1", &websocket.Conn{})
		client2 := NewClient("Room1", &websocket.Conn{})
		//wait for hub to finish register/unregister client first before continue
		var wg sync.WaitGroup

		wg.Add(2)
		h.Register <- client1 //client 1 join room 1
		go func() {
			for {
				if len(h.Rooms["Room1"].users) == 1 {
					break
				}
			}
			defer wg.Done()
		}()
		h.Register <- client2 //client 2 join room 1
		go func() {
			for {
				if len(h.Rooms["Room1"].users) == 2 {
					break
				}
			}
			defer wg.Done()
		}()
		wg.Wait()
		assert.Equal(t, 2, len(h.Rooms["Room1"].users))

		wg.Add(1)
		h.Unregister <- client1 //client 1 leave room 1
		go func() {
			for {
				if len(h.Rooms["Room1"].users) == 1 {
					break
				}
			}
			defer wg.Done()
		}()
		wg.Wait()
		assert.Equal(t, 1, len(h.Rooms["Room1"].users))
	})

	// t.Run("Server receives and broadcasts message correctly", func(t *testing.T) {
	// 	ctx, cancel := context.WithCancel(context.Background())
	// 	h := NewHub()
	// 	go h.Init(ctx)
	// 	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		upgrader := websocket.Upgrader{}
	// 		conn, err := upgrader.Upgrade(w, r, nil)
	// 		assert.NoError(t, err)

	// 		client := NewClient("Test Room", conn)
	// 		h.Register <- client
	// 		go h.ReadMsgFrom(client)
	// 	}))
	// 	defer server.Close()
	// 	defer cancel()
	// 	wsURL := "ws" + server.URL[4:]

	// 	user1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	// 	assert.NoError(t, err)

	// 	expected := []byte("Test Message")
	// 	err = user1.WriteMessage(websocket.TextMessage, expected)
	// 	assert.NoError(t, err)

	// 	assert.Equal(t, expected, <-h.Broadcast)
	// })

}

// func TestMonitorUser(t *testing.T) {
// 	t.Run("Server should terminate the chatroom when all user left", func(t *testing.T) {
// 		ctx, cancel := context.WithCancel(context.Background())
// 		room := New("Test Room", ctx)
// 		go room.Init()
// 		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 			upgrader := websocket.Upgrader{}
// 			conn, err := upgrader.Upgrade(w, r, nil)
// 			assert.NoError(t, err)
// 			room.Register <- conn
// 			go room.ReadMsgFrom(conn)
// 			go room.MonitorUser(cancel)
// 			err = room.ReadMsgFrom(conn)
// 			if err != nil {
// 				if websocket.IsCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 					log.Printf("connection from %s is closed", conn.RemoteAddr().String())
// 					room.Unregister <- conn //notify room when user left
// 				}
// 			}

// 		}))
// 		defer server.Close()

// 		wsURL := "ws" + server.URL[4:]

// 		user1, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 		assert.NoError(t, err)

// 		user2, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
// 		assert.NoError(t, err)

// 		user1.Close()
// 		user2.Close()

// 		// assert.
// 	})
// }
