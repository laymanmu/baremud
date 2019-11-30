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
	defer (Track("Start", g.log))()
	go g.handleNewClients()
	go g.handleClosedClients()
	go g.handleClientInput()
	g.server.Start()
	g.run()
}

// run is the run loop:
func (g *Game) run() {
	defer (Track("run", g.log))()
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
	defer (Track("handleNewClients", g.log))()
	for {
		client := <-g.newClients
		player := NewPlayer("player", client)
		g.players[player.ID] = player
		g.broadcast("[all] %s joined", player.ID)
		g.log("added %s for %s", player.ID, client.ID)
	}
}

// handleClosedClients removes closed clients
func (g *Game) handleClosedClients() {
	defer (Track("handleClosedClients", g.log))()
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

// handleClientInput handles client input from all the clients
func (g *Game) handleClientInput() {
	defer (Track("handleClientInput", g.log))()
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
	defer (Track("removePlayer", g.log))()
	delete(g.players, player.ID)
	g.broadcast("[all] %s left", player.Name)
	g.log("removed %s for %s", player.ID, player.client.ID)
}

// tick handles a single tick
func (g *Game) tick(tickCount int) {
	defer (Track(fmt.Sprintf("tick:%v", tickCount), g.log))()
	for _, player := range g.players {
		player.Update(g)
	}
}

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	defer (Track("broadcast", g.log))()
	msg := fmt.Sprintf(message, args...)
	for _, player := range g.players {
		player.client.Write(msg)
	}
}

// getPlayer gets a player from a client lookup
func (g *Game) getPlayer(client *Client) *Player {
	defer (Track("getPlayer", g.log))()
	for _, player := range g.players {
		if player.client.ID == client.ID {
			return player
		}
	}
	return nil
}
