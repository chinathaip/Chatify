package hub

import (
	"context"
	"encoding/json"
	"log"
	"net"
	"sync"
	"testing"

	"github.com/chinathaip/chatify/service"
	"github.com/google/uuid"
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

func syncRoomSize(h *H, room string, expectedLength int, wg *sync.WaitGroup) {
	for {
		if len(h.getRoom(room).users) == expectedLength {
			break
		}
	}
	defer wg.Done()
}

func TestInit(t *testing.T) {

	t.Run("Register new client should create new room if not exists", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		client := NewClient("Test room", &websocket.Conn{})

		h.Register <- client
		_, found := h.Rooms["Test room"]

		assert.Equal(t, 1, len(h.Rooms))
		assert.True(t, found)
	})

	t.Run("Register new client should add new user to the existing room", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		client2 := NewClient("Room1", &websocket.Conn{})
		h.Rooms["Room1"] = NewRoom("Room1")

		h.Register <- client2
		_, found1 := h.Rooms["Room1"]

		assert.Equal(t, 1, len(h.Rooms))
		assert.True(t, found1)
	})

	t.Run("Unregister should remove user from the existing room", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		client1 := NewClient("Room1", &websocket.Conn{})
		client2 := NewClient("Room1", &websocket.Conn{})
		r := NewRoom("Room1")
		r.setNewUser(client1.conn)
		r.setNewUser(client2.conn)
		h.setNewRoom(r.name, r)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		h.Unregister <- client1 //client 1 leave room 1
		go syncRoomSize(h, "Room1", 1, wg)

		wg.Wait()
		assert.Equal(t, 1, len(h.getRoom("Room1").users))
	})

	t.Run("Room should be terminated when last user left", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		client := NewClient("Room1", &websocket.Conn{})
		r := NewRoom("Room1")
		r.setNewUser(client.conn)
		h.setNewRoom(r.name, r)

		h.Unregister <- client

		assert.Nil(t, h.getRoom("Room1"))
	})

	t.Run("Broadcast message only within the room", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		mockConnection1 := newMock()
		mockConnection2 := newMock()
		mockConnection3 := newMock()
		r1 := NewRoom("Room1")
		r999 := NewRoom("Room999")
		r1.setNewUser(mockConnection1)
		r1.setNewUser(mockConnection2)
		r999.setNewUser(mockConnection3) //third client is assigned to different room
		h.setNewRoom(r1.name, r1)
		h.setNewRoom(r999.name, r999)

		wg := &sync.WaitGroup{}
		wg.Add(1)
		h.Broadcast <- &Message{"Room1", &JSONMessage{Sender: service.User{}, Text: "Hello"}} //message for user in room 1 only
		go func() {
			defer wg.Done()
			for {
				if mockConnection1.hasReceived && mockConnection2.hasReceived {
					break
				}
			}
		}()

		wg.Wait()
		assert.True(t, mockConnection1.hasReceived)
		assert.True(t, mockConnection2.hasReceived)
		assert.False(t, mockConnection3.hasReceived) //client 3 should not receive the message
		assert.Equal(t, 2, len(r1.users))
		assert.Equal(t, 1, len(r999.users))
		assert.Equal(t, 2, len(h.Rooms))
	})
}

func TestReadMsgFrom(t *testing.T) {
	t.Run("Read message from client correctly", func(t *testing.T) {
		h := New(nil, nil, nil)
		ctx, cancel := context.WithCancel(context.Background())
		go h.Init(ctx)
		defer cancel()
		mockConnection1 := newMock()
		mockConnection2 := newMock()
		client1 := NewClient("Room1", mockConnection1)
		client2 := NewClient("Room1", mockConnection2)
		r := NewRoom("Room1")
		r.setNewUser(client1.conn)
		r.setNewUser(client2.conn)
		h.setNewRoom(r.name, r)
		go h.ReadMsgFrom(client1)
		wg := &sync.WaitGroup{}

		id, _ := uuid.Parse("dapdkadakpk")
		msg := &JSONMessage{Type: "message", Sender: service.User{ID: id, Username: "Khing"}, Text: "Hello World"}
		data, err := json.Marshal(*msg)
		assert.NoError(t, err)
		log.Println("Here is the marshalled Data: ", data)

		wg.Add(1)
		err = client1.conn.WriteMessage(websocket.TextMessage, data)
		assert.NoError(t, err)
		result := <-h.Broadcast
		go func() {
			for {
				if mockConnection1.hasSent && mockConnection2.hasReceived {
					break
				}
			}
			defer wg.Done()
		}()

		wg.Wait()
		assert.Equal(t, "Hello World", result.data.Text)
		assert.True(t, mockConnection1.hasSent)
		assert.True(t, mockConnection2.hasReceived)
	})
}
