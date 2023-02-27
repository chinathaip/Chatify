package chatroom

type R struct {
	name  string
	users map[connection]bool
}

func New(roomName string) *R {
	return &R{
		name:  roomName,
		users: make(map[connection]bool),
	}
}
