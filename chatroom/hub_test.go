package chatroom

import (
	"context"
	"net"
	"sync"
	"testing"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
)

type mockConnection struct {
	msg         string
	hasReceived bool
	hasSent     bool
}

func newMock() *mockConnection {
	return &mockConnection{hasReceived: false, hasSent: false}
}

func (c *mockConnection) RemoteAddr() net.Addr {
	return &net.IPAddr{}
}

func (c *mockConnection) WriteMessage(messageType int, data []byte) error {
	c.msg = string(data)
	c.hasReceived = true
	return nil
}

func (c *mockConnection) ReadMessage() (messageType int, data []byte, err error) {
	c.hasSent = true
	return 0, []byte(c.msg), nil
}

func syncRoomSize(h *Hub, room string, expectedLength int, wg *sync.WaitGroup) {
	for {
		if len(h.Rooms[room].users) == expectedLength {
			break
		}
	}
	defer wg.Done()
}

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
		wg := &sync.WaitGroup{}
		wg.Add(2)
		h.Register <- client1 //client 1 join room 1
		go syncRoomSize(h, "Room1", 1, wg)
		h.Register <- client2 //client 2 join room 1
		go syncRoomSize(h, "Room1", 2, wg)
		wg.Wait() //wait for both register to finish
		assert.Equal(t, 2, len(h.Rooms["Room1"].users))

		wg.Add(1)
		h.Unregister <- client1 //client 1 leave room 1
		syncRoomSize(h, "Room1", 1, wg)

		wg.Wait()
		assert.Equal(t, 1, len(h.Rooms["Room1"].users))
	})

	t.Run("Room should be terminated when last user left", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		client := NewClient("Room1", &websocket.Conn{})
		wg := &sync.WaitGroup{}
		wg.Add(1)
		h.Register <- client
		go syncRoomSize(h, "Room1", 1, wg)
		wg.Wait()
		assert.Equal(t, 1, len(h.Rooms["Room1"].users))

		wg.Add(1)
		h.Unregister <- client
		go syncRoomSize(h, "Room1", 0, wg)

		wg.Wait()
		assert.Nil(t, h.Rooms["Room1"])
	})

	t.Run("Broadcast message only within the room", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		mockConnection1 := newMock()
		mockConnection2 := newMock()
		mockConnection3 := newMock()
		client1 := NewClient("Room1", mockConnection1)
		client2 := NewClient("Room1", mockConnection2)
		client3 := NewClient("Room999", mockConnection3) //different room
		wg := &sync.WaitGroup{}
		wg.Add(3)
		h.Register <- client1
		go syncRoomSize(h, "Room1", 1, wg)
		h.Register <- client2
		go syncRoomSize(h, "Room1", 2, wg)
		h.Register <- client3
		go syncRoomSize(h, "Room999", 1, wg)
		wg.Wait()
		assert.Equal(t, 2, len(h.Rooms["Room1"].users))
		assert.Equal(t, 1, len(h.Rooms["Room999"].users))

		wg.Add(1)
		h.Broadcast <- &Message{"Room1", []byte("Hello!")} //message for user in room 1 only
		go func() {
			for {
				if mockConnection1.hasReceived && mockConnection2.hasReceived {
					break
				}
			}
			defer wg.Done()
		}()

		wg.Wait()
		assert.True(t, mockConnection1.hasReceived)
		assert.True(t, mockConnection2.hasReceived)
		assert.False(t, mockConnection3.hasReceived) //client 3 should not receive the message
	})
}

func TestReadMsgFrom(t *testing.T) {
	t.Run("Read message from client correctly", func(t *testing.T) {
		h := NewHub()
		go h.Init(context.Background())
		mockConnection1 := newMock()
		mockConnection2 := newMock()
		client1 := NewClient("Room1", mockConnection1)
		client2 := NewClient("Room1", mockConnection2)
		wg := &sync.WaitGroup{}
		wg.Add(2)
		h.Register <- client1
		go h.ReadMsgFrom(client1)
		go syncRoomSize(h, "Room1", 1, wg)
		h.Register <- client2
		go syncRoomSize(h, "Room1", 2, wg)
		wg.Wait()
		assert.Equal(t, 2, len(h.Rooms["Room1"].users))

		wg.Add(1)
		err := client1.conn.WriteMessage(websocket.TextMessage, []byte("Hello World"))
		assert.NoError(t, err)

		result := <-h.Broadcast
		assert.Equal(t, "Hello World", string(result.data))
		go func() {
			for {
				if mockConnection1.hasSent && mockConnection2.hasReceived {
					break
				}
			}
			defer wg.Done()
		}()
		wg.Wait()
		assert.True(t, mockConnection1.hasSent)
		assert.True(t, mockConnection2.hasReceived)
	})
}
