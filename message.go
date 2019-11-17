package main

import (
	"fmt"

	"github.com/google/uuid"
)

// MessageType is a type of message
type MessageType int

// message types:
const (
	ErrorMessage         MessageType = iota
	ClientStartedMessage MessageType = iota
	ClientStoppedMessage MessageType = iota
	ClientChatMessage    MessageType = iota
	ClientLookMessage    MessageType = iota
	ClientEnterMessage   MessageType = iota
)

// PrintMessageTypeValues will print message type values
func PrintMessageTypeValues() {
	fmt.Printf("MessageType values:\n")
	fmt.Printf("  ErrorMessage:         %v\n", ErrorMessage)
	fmt.Printf("  ClientStartedMessage: %v\n", ClientStartedMessage)
	fmt.Printf("  ClientStoppedMessage: %v\n", ClientStoppedMessage)
	fmt.Printf("  ClientChatMessage:    %v\n", ClientChatMessage)
	fmt.Printf("  ClientLookMessage:    %v\n", ClientLookMessage)
	fmt.Printf("  ClientEnterMessage:   %v\n", ClientEnterMessage)
}

// GetMessageTypeName returns a message type name
func GetMessageTypeName(msgType MessageType) string {
	switch msgType {
	case ErrorMessage:
		return "ErrorMessage"
	case ClientStartedMessage:
		return "ClientStartedMessage"
	case ClientStoppedMessage:
		return "ClientStoppedMessage"
	case ClientChatMessage:
		return "ClientChatMessage"
	case ClientLookMessage:
		return "ClientLookMessage"
	case ClientEnterMessage:
		return "ClientEnterMessage"
	default:
		return fmt.Sprintf("Unknown MessageType: %v", msgType)
	}
}

// Message is a message
type Message struct {
	ID      string
	Type    MessageType
	Name    string
	Client  *Client
	Message string
	Args    []string
}

// NewMessage creates a message
func NewMessage(msgType MessageType, client *Client, message string, args []string) *Message {
	ID := uuid.New().String()
	if args == nil {
		args = []string{}
	}
	Name := GetMessageTypeName(msgType)
	return &Message{ID, msgType, Name, client, message, args}
}

// NewErrorMessage creates a message
func NewErrorMessage(client *Client, message string) *Message {
	return NewMessage(ErrorMessage, client, message, nil)
}

// NewClientStartedMessage creates a message
func NewClientStartedMessage(client *Client) *Message {
	return NewMessage(ClientStartedMessage, client, "", nil)
}

// NewClientStoppedMessage creates a message
func NewClientStoppedMessage(client *Client) *Message {
	return NewMessage(ClientStoppedMessage, client, "", nil)
}

// NewClientLookMessage creates a message
func NewClientLookMessage(client *Client, args []string) *Message {
	return NewMessage(ClientLookMessage, client, "", args)
}

// NewClientEnterMessage creates a message
func NewClientEnterMessage(client *Client, gateName string) *Message {
	return NewMessage(ClientEnterMessage, client, gateName, nil)
}

// NewClientChatMessage creates a message
func NewClientChatMessage(client *Client, message string) *Message {
	return NewMessage(ClientChatMessage, client, message, nil)
}
