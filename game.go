package main

import (
	"fmt"
	"log"
	"time"
)

// Game is a game
type Game struct {
	ID          string
	server      *Server
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
	return &Game{NewID(), server, clientInput, newClients, clients}
}

// Start starts the game
func (g *Game) Start() {
	g.server.Start()
	go g.handleNewClients()
	go g.handleClientInput()
	g.run()
}

// run is the run loop:
func (g *Game) run() {
	tickTime := time.Duration(5000) * time.Millisecond
	tickCount := 0
	for {
		select {
		case <-time.After(tickTime):
			tickCount++
			g.tick(tickCount)
			g.pruneClients()
		}
	}
}

// tick handles a single tick
func (g *Game) tick(tickCount int) {
	msg := fmt.Sprintf("tick: %v", tickCount)
	g.log(msg)
}

// handleNewClients handles the new clients channel
func (g *Game) handleNewClients() {
	for {
		client := <-g.newClients
		g.clients[client.ID] = client
		g.broadcast("[all] %s joined", client.ID)
	}
}

// handleClientInput handles client input from all the clients
func (g *Game) handleClientInput() {
	for {
		msg := <-g.clientInput
		g.broadcast("[all] %s says: %s", msg.Client.ID, msg.Input)
	}
}

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	for _, client := range g.clients {
		client.Write(msg)
	}
}

// removeClient removes a client:
func (g *Game) removeClient(client *Client) {
	delete(g.clients, client.ID)
	g.broadcast("[all] %s left", client.ID)
}

// pruneClients removes closed clients
func (g *Game) pruneClients() {
	for _, client := range g.clients {
		if client.IsClosed {
			g.removeClient(client)
		}
	}
}

// log is for logging a message
func (g *Game) log(message string, args ...interface{}) {
	msg := fmt.Sprintf("game:%s | %s\n", g.ID, message)
	log.Printf(msg, args...)
}
