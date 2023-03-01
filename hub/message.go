package hub

type JSONMessage struct {
	Type     string `json:"type"`
	SenderID string `json:"sender_id"`
	Text     string `json:"text"`
}
