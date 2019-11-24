package main

import (
	"bytes"
	"fmt"
	"strings"
)

// Commander handles commands
type Commander struct {
	commands map[string]*Command
}

// NewCommander creates a commander
func NewCommander() *Commander {
	c := make(map[string]*Command)
	c["help"] = NewCommand("help", "list all commands or get details about a single comand", 0, 0)
	c["look"] = NewCommand("look", "look at the room or something else", 0, 1)
	c["enter"] = NewCommand("enter", "enter a gate to change rooms", 1, 0)
	c["exit"] = NewCommand("exit", "exit the game", 0, 0)
	c["say"] = NewCommand("say", "send a chat message", 1, 0)
	c["mkroom"] = NewCommand("mkroom", "make a new room", 0, 0)
	return &Commander{commands: c}
}

// ListCommands returns a list of all commands
func (c *Commander) ListCommands() string {
	var b bytes.Buffer
	for command := range c.commands {
		b.WriteString(fmt.Sprintf("%s - %s\r\n", command, c.commands[command].Desc))
	}
	return b.String()
}

// IsCommand will check if a given string is a known command
func (c *Commander) IsCommand(command string) bool {
	_, found := c.commands[command]
	return found
}

// HandleHelp handles a command message
func (c *Commander) HandleHelp(msg *CommandMessage) {
	command, ok := c.commands[msg.Arguments[0]]
	if ok {
		msg.Meta.Client.Write(command.Help())
	} else {
		msg.Meta.Client.Write(c.ListCommands())
	}
}

// HandleLook describes a room
func (c *Commander) HandleLook(msg *CommandMessage) {
	msg.Meta.Client.Write(msg.Meta.Client.Player.Room.Look())
}

// HandleEnter moves a client through a gate
func (c *Commander) HandleEnter(msg *CommandMessage) {
	msg.Meta.Client.Player.EnterGate(msg.Arguments[0])
}

// HandleExit removes a player from the game
func (c *Commander) HandleExit(msg *CommandMessage) {
	// nothing to do
}

// HandleSay sends a message to the room
func (c *Commander) HandleSay(msg *CommandMessage) {
	message := fmt.Sprintf("%s says: %s", msg.Meta.Client.Player.Name, strings.Join(msg.Arguments, " "))
	msg.Meta.Client.Player.Room.Broadcast(message)
}

// HandleMkroom makes a room
func (c *Commander) HandleMkroom(msg *CommandMessage) {
	msg.Meta.Client.Player.Room.Broadcast("todo: make a room")
}
