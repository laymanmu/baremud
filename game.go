package main

import (
	"fmt"
	"time"
)

// Game is a game
type Game struct {
	ID          string
	log         Logger
	newClients  <-chan *Client
	clientInput <-chan *ClientInputMessage
	players     map[string]*Player
	server      *Server
	commander   *Commander
}

// NewGame creates a world
func NewGame() *Game {
	id := NewID("game")
	log := NewLogger(id)
	newClients := make(chan *Client)
	clientInput := make(chan *ClientInputMessage)
	players := make(map[string]*Player)
	server := NewServer(":2323", newClients, clientInput)
	commander := NewCommander()
	return &Game{id, log, newClients, clientInput, players, server, commander}
}

// Start starts the game
func (g *Game) Start() {
	g.log("Start started")
	defer func() { g.log("Start stopped") }()
	go g.handleNewClients()
	go g.handleClosedClients()
	go g.handleClientInput()
	g.server.Start()
	g.run()
}

// run is the run loop:
func (g *Game) run() {
	g.log("run started")
	defer func() { g.log("run stopped") }()
	tickTime := time.Duration(5000) * time.Millisecond
	tickCount := 0
	for {
		select {
		case <-time.After(tickTime):
			tickCount++
			g.tick(tickCount)
		}
	}
}

// handleNewClients adds new clients
func (g *Game) handleNewClients() {
	g.log("handleNewClients started")
	defer func() { g.log("handleNewClients stopped") }()
	for {
		client := <-g.newClients
		player := NewPlayer("player", client)
		g.players[player.ID] = player
		g.broadcast("[all] %s joined", player.ID)
		g.log("added player:%s", player.ID)
	}
}

// handleClosedClients removes closed clients
func (g *Game) handleClosedClients() {
	defer g.track("handleClosedClients")
	interval := time.Duration(1000) * time.Millisecond
	for {
		select {
		case <-time.After(interval):
			closed := []*Player{}
			for _, player := range g.players {
				if player.client.IsClosed {
					closed = append(closed, player)
				}
			}
			for _, player := range closed {
				g.removePlayer(player)
			}
		}
	}
}

func (g *Game) track(funcName string) func() {
	g.log("%s started", funcName)
	return func() {
		g.log("%s stopped", funcName)
	}
}

// handleClientInput handles client input from all the clients
func (g *Game) handleClientInput() {
	g.log("handleClientInput started")
	defer func() { g.log("handleClientInput stopped") }()
	for {
		msg := <-g.clientInput
		cmd, err := g.commander.GetCommand(msg.Input)
		player := g.getPlayer(msg.Client)
		if err != nil {
			msg.Client.Write(err.Error())
			continue
		}
		if cmd.Name == "exit" {
			g.removePlayer(player)
			msg.Client.Close()
			continue
		}
		go g.commander.HandleCommand(cmd, player, g)
	}
}

// removePlayer removes a player
func (g *Game) removePlayer(player *Player) {
	delete(g.players, player.ID)
	g.broadcast("[all] %s left", player.Name)
	g.log("removed client:%s", player.Name)
}

// tick handles a single tick
func (g *Game) tick(tickCount int) {
	g.log("tick: %v", tickCount)
	for _, player := range g.players {
		player.Update(g)
	}
}

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	for _, player := range g.players {
		player.client.Write(msg)
	}
}

// getPlayer gets a player from a client lookup
func (g *Game) getPlayer(client *Client) *Player {
	for _, player := range g.players {
		if player.client.ID == client.ID {
			return player
		}
	}
	return nil
}
