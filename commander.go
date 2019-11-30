package main

import (
	"fmt"
	"strings"
)

// Commander is a commander
type Commander struct {
	ID       string
	commands map[string]*Command
	log      Logger
}

// NewCommander creates a commander
func NewCommander() *Commander {
	id := NewID("commander")
	log := NewLogger(id)
	help := NewCommand("help", "shows help")
	look := NewCommand("look", "shows surroundings")
	enter := NewCommand("enter", "enters a gate")
	exit := NewCommand("exit", "enters a gate")
	say := NewCommand("say", "send a message to chat")
	stats := NewCommand("stats", "shows player stats")
	debug := NewCommand("debug", "debug")
	commands := make(map[string]*Command)
	for _, cmd := range []*Command{help, look, enter, exit, say, stats, debug} {
		commands[cmd.Name] = cmd
	}
	return &Commander{id, commands, log}
}

// GetCommand returns a command from client input if valid
func (c *Commander) GetCommand(input string) (*Command, error) {
	defer (Track("GetCommand", c.log))()
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
	defer (Track("CommandNames", c.log))()
	names := []string{}
	for name := range c.commands {
		names = append(names, name)
	}
	return names
}

// HandleCommand handles a command
func (c *Commander) HandleCommand(command *Command, player *Player, game *Game) {
	defer (Track("HandleCommand", c.log))()
	client := player.client
	switch command.Name {
	case "look":
		if command.HasArgs() {
			client.Write("you look at %s", command.ArgString())
		} else {
			client.Write("you look around")
		}
	case "enter":
		client.Write("you enter %s", command.Args[0])
	case "help":
		client.Write("commands: %v", c.CommandNames())
	case "say":
		game.broadcast("[all] %s: %s", client.ID, command.ArgString())
	case "stats":
		client.Write(player.BuildPrompt())
	case "debug":
		player.Resources["health"].Value = 0
		player.Resources["energy"].Value = 0
		game.log("len(game.players): %v", len(game.players))
	default:
		c.log("unhandled command: %s", command.Name)
	}
}
