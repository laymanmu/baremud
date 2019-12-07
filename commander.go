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
	if len(fields) == 0 {
		return nil, nil
	}
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
		player.Look()
	case "enter":
		client.Write("you enter %s", command.Args[0])
	case "help":
		client.Write("commands: %v", c.CommandNames())
	case "say":
		game.broadcast("[chat] %s: %s", player.ID, command.ArgString())
	case "stats":
		client.Write(player.BuildPrompt())
	case "debug":
		client.Write(player.Place.Look(player))
		for _, p := range game.players {
			client.Write(fmt.Sprintf("  %s %s", p.ID, p.BuildPrompt()))
		}
	default:
		msg := fmt.Sprintf("error: unhandled command: %s", command.Name)
		c.log(msg)
		client.Write(msg)
	}
}
