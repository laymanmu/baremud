package network

// Message is a network message
type Message struct {
	From    string
	Message string
}

// NewMessage creates a network message
func NewMessage(from, message string) *Message {
	return &Message{From: from, Message: message}
}
