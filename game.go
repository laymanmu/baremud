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
	state       *State
	server      *Server
	commander   *Commander
	players     []*Player
}

// NewGame creates a world
func NewGame() *Game {
	id := NewID("game")
	log := NewLogger(id)
	newClients := make(chan *Client)
	clientInput := make(chan *ClientInputMessage)
	state := NewState()
	server := NewServer(":2323", newClients, clientInput)
	commander := NewCommander()
	players := []*Player{}
	return &Game{id, log, newClients, clientInput, state, server, commander, players}
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
	for {
		select {
		case <-time.After(tickTime):
			g.state.TickCount++
			g.tick(g.state.TickCount)
		}
	}
}

// tick handles a single tick
func (g *Game) tick(tickCount int) {
	defer (Track(fmt.Sprintf("-= tick:%v =-", tickCount), g.log))()

	// update the players list:
	g.players = g.state.GetAllPlayers()

	// update each player:
	for _, player := range g.players {
		player.Update(g)
	}
}

// handleNewClients adds new clients
func (g *Game) handleNewClients() {
	n := "handleNewClients"
	defer (Track(n, g.log))()
	for {
		client := <-g.newClients
		player := NewPlayer("player", client)
		g.state.Add(player, g.state.StartingPlace())
		g.broadcast("%s joined", player.ID)
		g.log("%s | added %s for %s", n, player.ID, client.ID)
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
		if err != nil {
			msg.Client.Write(err.Error())
			continue
		}
		player := g.state.clientPlayer[msg.Client.ID]
		if cmd.Name == "exit" {
			g.removePlayer(player)
			msg.Client.Close()
			continue
		}
		go g.commander.HandleCommand(cmd, player, g, g.state)
	}
}

// removePlayer removes a player
func (g *Game) removePlayer(player *Player) {
	defer (Track("removePlayer", g.log))()
	g.state.RemovePlayer(player)
	g.broadcast("%s left", player.ID)
	g.log("removed %s for %s", player.ID, player.client.ID)
}

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	defer (Track("broadcast", g.log))()
	msg := fmt.Sprintf(message, args...)
	for _, player := range g.players {
		player.client.Write(msg)
	}
}
