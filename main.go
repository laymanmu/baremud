package main

import (
	"strings"

	"github.com/google/uuid"
)

// NewID creates a unique id
func NewID() string {
	id := uuid.New().String()
	i := strings.IndexByte(id, '-')
	if i == -1 {
		return id
	}
	return id[:i]
}

// ClientInputMessage is a message with input from a client
type ClientInputMessage struct {
	ID     string
	Client *Client
	Input  string
}

// NewClientInputMessage creates a message
func NewClientInputMessage(client *Client, input string) *ClientInputMessage {
	return &ClientInputMessage{NewID(), client, input}
}

func main() {
	game := NewGame()
	game.Start()
}
