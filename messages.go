package main

// MessageMetaData holds the message identifiers shared by all message types
type MessageMetaData struct {
	ID     string
	Name   string
	Client *Client
}

// message types:

// ErrorMessage is a message for an error
type ErrorMessage struct {
	Meta    *MessageMetaData
	Message string
}

// ClientJoinMessage is a message for requesting a client to join
type ClientJoinMessage struct {
	Meta *MessageMetaData
}

// ClientLeaveMessage is a message for requesting a client to leave
type ClientLeaveMessage struct {
	Meta   *MessageMetaData
	Reason string
}

// CommandMessage is a message for requesting a command to be run
type CommandMessage struct {
	Meta      *MessageMetaData
	Command   string
	Arguments []string
}

// InputMessage is a message with new input from a client
type InputMessage struct {
	Meta  *MessageMetaData
	Input string
}

// OutputMessage is a message with output for a client
type OutputMessage struct {
	Meta   *MessageMetaData
	Output string
}

// constructors:

// NewMessageMetaData creates a MessageMetaData
func NewMessageMetaData(name string, client *Client) *MessageMetaData {
	return &MessageMetaData{GetID(), name, client}
}

// NewClientJoinMessage creates a message for requesting a client to join
func NewClientJoinMessage(client *Client) *ClientJoinMessage {
	meta := NewMessageMetaData("ClientJoinMessage", client)
	return &ClientJoinMessage{meta}
}

// NewClientLeaveMessage creates a message for requesting a client to leave
func NewClientLeaveMessage(client *Client, reason string) *ClientLeaveMessage {
	meta := NewMessageMetaData("ClientLeaveMessage", client)
	return &ClientLeaveMessage{meta, reason}
}

// NewCommandMessage creates a message for requesting a command to be run
func NewCommandMessage(client *Client, command string, args []string) *CommandMessage {
	meta := NewMessageMetaData("CommandMessage", client)
	return &CommandMessage{meta, command, args}
}

// NewErrorMessage creates a message for an error
func NewErrorMessage(client *Client, message string) *ErrorMessage {
	meta := NewMessageMetaData("ErrorMessage", client)
	return &ErrorMessage{meta, message}
}

// NewInputMessage creates a message with new input from the client
func NewInputMessage(client *Client, input string) *InputMessage {
	meta := NewMessageMetaData("InputMessage", client)
	return &InputMessage{meta, input}
}

// NewOutputMessage creates a message with new input from the client
func NewOutputMessage(client *Client, output string) *OutputMessage {
	meta := NewMessageMetaData("OutputMessage", client)
	return &OutputMessage{meta, output}
}
