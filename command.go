package main

import (
	"bytes"
	"fmt"
)

// Command is a command
type Command struct {
	Name            string
	Desc            string
	NumRequiredArgs int
	NumOptionalArgs int
	Examples        []string
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
		for _, example := range c.Examples {
			b.WriteString(fmt.Sprintf("  %s\r\n", example))
		}
	}
	return b.String()
}
