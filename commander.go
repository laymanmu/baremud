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
	commands := make(map[string]*Command)
	commander := &Commander{id, commands, log}

	return commander
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
	switch command.Name {
	case "look":
		c.handleLook(command, player, game)
	case "enter":
		c.handleEnter(command, player, game)
	case "help":
		c.handleHelp(command, player, game)
	case "say":
		c.handleSay(command, player, game)
	case "stats":
		c.handleStats(command, player, game)
	case "debug":
		c.handleDebug(command, player, game)
	default:
		msg := fmt.Sprintf("error: unhandled command: %s", command.Name)
		c.log(msg)
		player.client.Write(msg)
	}
}

// createCommands will create the commands
func (c *Commander) createCommands() {
	c.commands["help"] = NewCommand("help", "shows help")
	c.commands["look"] = NewCommand("look", "shows surroundings")
	c.commands["enter"] = NewCommand("enter", "enters a gate")
	c.commands["exit"] = NewCommand("exit", "enters a gate")
	c.commands["say"] = NewCommand("say", "send a message to chat")
	c.commands["stats"] = NewCommand("stats", "shows player stats")
	c.commands["debug"] = NewCommand("debug", "debug")
}

// handleLook handles command: look
func (c *Commander) handleLook(command *Command, player *Player, game *Game) {
	player.Look()
}

// handleEnter handles command: enter
func (c *Commander) handleEnter(command *Command, player *Player, game *Game) {
	player.client.Write("you enter %s", command.Args[0])
}

// handleHelp handles command: help
func (c *Commander) handleHelp(command *Command, player *Player, game *Game) {
	player.client.Write("commands: %v", c.CommandNames())
}

// handleSay handles command: say
func (c *Commander) handleSay(command *Command, player *Player, game *Game) {
	game.broadcast("[chat] %s: %s", player.ID, command.ArgString())
}

// handleStats handles command: stats
func (c *Commander) handleStats(command *Command, player *Player, game *Game) {
	player.client.Write(player.BuildPrompt())
}

// handleDebug handles command: debug
func (c *Commander) handleDebug(command *Command, player *Player, game *Game) {
	player.client.Write(player.Place.Look(player))
	for _, p := range game.players {
		player.client.Write(fmt.Sprintf("  %s %s", p.ID, p.BuildPrompt()))
	}
}
