package main

import (
	"fmt"
	"time"
)

// Game is a game
type Game struct {
	ID          string
	server      *Server
	commander   *Commander
	clientInput <-chan *ClientInputMessage
	newClients  <-chan *Client
	clients     map[string]*Client
}

// NewGame creates a world
func NewGame() *Game {
	clients := make(map[string]*Client)
	newClients := make(chan *Client)
	clientInput := make(chan *ClientInputMessage)
	server := NewServer(":2323", newClients, clientInput)
	commander := NewCommander()
	return &Game{NewID(), server, commander, clientInput, newClients, clients}
}

// Start starts the game
func (g *Game) Start() {
	g.log("Start started")
	defer func() { g.log("Start stopped") }()
	g.server.Start()
	go g.handleNewClients()
	go g.handleClientInput()
	go g.handleClientPruning()
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
		g.clients[client.ID] = client
		g.broadcast("[all] %s joined", client.ID)
		g.log("added client:%s", client.ID)
	}
}

// handleClientPruning removes closed clients
func (g *Game) handleClientPruning() {
	g.log("handleClientPruning started")
	defer func() { g.log("handleClientPruning stopped") }()
	interval := time.Duration(1000) * time.Millisecond
	for {
		select {
		case <-time.After(interval):
			for _, client := range g.clients {
				if client.IsClosed {
					g.removeClient(client)
				}
			}
		}
	}
}

// handleClientInput handles client input from all the clients
func (g *Game) handleClientInput() {
	g.log("handleClientInput started")
	defer func() { g.log("handleClientInput stopped") }()
	for {
		msg := <-g.clientInput
		cmd, err := g.commander.GetCommand(msg.Input)
		if err != nil {
			msg.Client.Write(err.Error())
			continue
		}
		if cmd.Name == "exit" {
			g.removeClient(msg.Client)
			msg.Client.Close()
			continue
		}
		go g.commander.HandleCommand(cmd, msg.Client)
	}
}

// removeClient removes a client:
func (g *Game) removeClient(client *Client) {
	delete(g.clients, client.ID)
	g.broadcast("[all] %s left", client.ID)
	g.log("removed client:%s", client.ID)
}

// tick handles a single tick
func (g *Game) tick(tickCount int) {
	g.log("tick: %v", tickCount)
}

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	for _, client := range g.clients {
		client.Write(msg)
	}
}

// log is for logging a message
func (g *Game) log(message string, args ...interface{}) {
	src := fmt.Sprintf("game:%s", g.ID)
	msg := fmt.Sprintf(message, args...)
	Log(src, msg)
}
