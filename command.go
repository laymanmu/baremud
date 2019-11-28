package main

import (
	"fmt"
	"strings"
)

// Command defines a command
type Command struct {
	Name string
	Desc string
	Args []string
}

// NewCommand creates a command
func NewCommand(name, desc string) *Command {
	args := []string{}
	return &Command{name, desc, args}
}

// Commander is a commander
type Commander struct {
	commands map[string]*Command
}

// NewCommander creates a commander
func NewCommander() *Commander {
	help := NewCommand("help", "shows help")
	look := NewCommand("look", "shows surroundings")
	enter := NewCommand("enter", "enters a gate")
	exit := NewCommand("exit", "enters a gate")
	debug := NewCommand("debug", "debug")
	commands := make(map[string]*Command)
	for _, cmd := range []*Command{help, look, enter, exit, debug} {
		commands[cmd.Name] = cmd
	}
	return &Commander{commands}
}

// GetCommand returns a command from client input if valid
func (c *Commander) GetCommand(input string) (*Command, error) {
	if len(input) < 1 {
		return nil, fmt.Errorf("missing input")
	}
	fields := strings.Fields(input)
	if command, ok := c.commands[fields[0]]; ok {
		if len(fields) > 1 {
			command.Args = fields[1:]
		}
		return command, nil
	}
	return nil, fmt.Errorf("command not found: %s. try: %v", fields[0], c.CommandNames())
}

// CommandNames returns a list of known command names
func (c *Commander) CommandNames() []string {
	names := []string{}
	for name := range c.commands {
		names = append(names, name)
	}
	return names
}

// HandleCommand handles a command
func (c *Commander) HandleCommand(command *Command, client *Client) {
	switch command.Name {
	case "look":
		client.Write("you look around")
	case "enter":
		client.Write("you enter %s", command.Args[0])
	case "help":
		client.Write("commands: %v", c.CommandNames())
	case "debug":
		client.Write("client.IsClosed:%v", client.IsClosed)
	default:
		client.Write("todo: handle command: %s", command.Name)
	}
}
