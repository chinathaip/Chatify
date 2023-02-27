package hub

type Room struct {
	name  string
	users map[connection]bool
}

func NewRoom(roomName string) *Room {
	return &Room{
		name:  roomName,
		users: make(map[connection]bool),
	}
}
