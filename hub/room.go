package hub

import "sync"

type Room struct {
	id    int
	name  string
	users map[connection]bool
	mutex sync.RWMutex
}

func NewRoom(roomName string) *Room {
	return &Room{
		name:  roomName,
		users: make(map[connection]bool),
		mutex: sync.RWMutex{},
	}
}

func (r *Room) setNewUser(c connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.users[c] = true
}

func (r *Room) deleteUser(c connection) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	delete(r.users, c)
}
