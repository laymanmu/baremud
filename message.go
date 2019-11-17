package main

import (
	"fmt"
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
	Name := GetMessageTypeName(msgType)
	return &Message{GetID(), msgType, Name, client, message, args}
}

// String will return a string representation which is good for logging
func (m *Message) String() string {
	format := "%s client:%s id:%s message:\"%s\" args:%s"
	return fmt.Sprintf(format, m.Name, m.Client.ID, m.ID, m.Message, m.Args)
}
