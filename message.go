package main

// ClientInputMessage is a message with input from a client
type ClientInputMessage struct {
	ID     string
	Client *Client
	Input  string
}

// NewClientInputMessage creates a message
func NewClientInputMessage(client *Client, input string) *ClientInputMessage {
	return &ClientInputMessage{NewID("message"), client, input}
}
