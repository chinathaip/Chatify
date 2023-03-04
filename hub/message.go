package hub

import "github.com/chinathaip/chatify/service"

type JSONMessage struct {
	Type   string       `json:"type"`
	Sender service.User `json:"sender"`
	Text   string       `json:"text"`
}
