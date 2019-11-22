package main

import (
	"bytes"
	"fmt"
)

// Commander handles commands
type Commander struct {
	commands map[string]*Command
}

// NewCommander creates a commander
func NewCommander() *Commander {
	c := make(map[string]*Command)
	c["help"]   = NewCommand("help",   "list all commands or get details about a single comand", 0, 0)
	c["look"]   = NewCommand("look",   "look at the room or something else", 0, 1)
	c["enter"]  = NewCommand("enter",  "enter a gate to change rooms", 1, 0)
	c["exit"]   = NewCommand("exit",   "exit the game", 0, 0)
	c["mkroom"] = NewCommand("mkroom", "make a new room", 0, 0)
	return &Commander{commands:c}
}

// AddCommand adds a command
func (c *Commander) AddCommand(command *Command) {
	c.commands[command.Name] = command
}

// ListCommands returns a list of all commands
func (c *Commander) ListCommands() string {
	var b bytes.Buffer
	for command, _ := range c.commands {
		b.WriteString(fmt.Sprintf("%s - %s\r\n", command, c.commands[command].Desc))
	}
	return b.String()
}

// Help shows help for a given command
func (c *Commander) Help(command string) string {
	com, ok := c.commands[command]; if ok {
		return com.Help()
	}
	return c.ListCommands()
}


// Command is a command 
type Command struct {
	Name string
	Desc string
	NumRequiredArgs int
	NumOptionalArgs int
	Examples []string
}

// NewCommand creates a command
func NewCommand(name, desc string, numRequiredArgs int, numOptionalArgs int) *Command {
	examples := []string{}
	return &Command{name, desc, numRequiredArgs, numOptionalArgs, examples}
}

// AddExample adds an example to show in help
func (c *Command) AddExample(example string) {
	c.Examples = append(c.Examples, example)
}

// Help returns the command help documentation
func (c *Command) Help() string {
	var b bytes.Buffer
	b.WriteString(fmt.Sprintf("%s - %s\r\n", c.Name, c.Desc))
	if len(c.Examples) > 0 {
		b.WriteString("Examples:\r\n")
		for example := range c.Examples {
			b.WriteString(fmt.Sprintf("  %s\r\n", example))
		}
	}
	return b.String()
}