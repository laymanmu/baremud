package main

import "github.com/google/uuid"

// ErrorMessage is a message
type ErrorMessage struct {
	ID      string
	Client  *Client
	Message string
}

// ClientLoggedOnMessage is a message
type ClientLoggedOnMessage struct {
	ID     string
	Client *Client
}

// ClientClosedMessage is a message
type ClientClosedMessage struct {
	ID     string
	Client *Client
}

// ClientLookMessage is a message
type ClientLookMessage struct {
	ID     string
	Client *Client
	Args   []string
}

// ClientChatMessage is a message
type ClientChatMessage struct {
	ID      string
	Client  *Client
	Message string
}

// ClientEnterMessage is a message
type ClientEnterMessage struct {
	ID     string
	Client *Client
	Args   []string
}

// NewErrorMessage creates a message
func NewErrorMessage(client *Client, message string) *ErrorMessage {
	return &ErrorMessage{ID: uuid.New().String(), Client: client, Message: message}
}

// NewClientLoggedOnMessage creates a message
func NewClientLoggedOnMessage(client *Client) *ClientLoggedOnMessage {
	return &ClientLoggedOnMessage{ID: uuid.New().String(), Client: client}
}

// NewClientClosedMessage creates a message
func NewClientClosedMessage(client *Client) *ClientClosedMessage {
	return &ClientClosedMessage{ID: uuid.New().String(), Client: client}
}

// NewClientLookMessage creates a message
func NewClientLookMessage(client *Client, args []string) *ClientLookMessage {
	return &ClientLookMessage{ID: uuid.New().String(), Client: client, Args: args}
}

// NewClientEnterMessage creates a message
func NewClientEnterMessage(client *Client, args []string) *ClientEnterMessage {
	return &ClientEnterMessage{ID: uuid.New().String(), Client: client, Args: args}
}

// NewClientChatMessage creates a message
func NewClientChatMessage(client *Client, message string) *ClientChatMessage {
	return &ClientChatMessage{ID: uuid.New().String(), Client: client, Message: message}
}
