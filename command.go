package main

import (
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

// ArgString returns the Args array as a single space delimited string
func (c *Command) ArgString() string {
	return strings.Join(c.Args, " ")
}

// FirstArg returns the first element of the Args array
func (c *Command) FirstArg() string {
	if c.HasArgs() {
		return c.Args[0]
	}
	return ""
}

// HasArgs returns a bool to specify if the command has anything in the Args array
func (c *Command) HasArgs() bool {
	return len(c.Args) > 0
}
