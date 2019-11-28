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
	g.server.Start()
	go g.handleNewClients()
	go g.handleClientInput()
	go g.handleClientPruning()
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
		cmd, err := g.commander.GetCommand(msg.Input)
		if err != nil {
			msg.Client.Write(err.Error())
			continue
		}
		if cmd.Name == "exit" {
			g.log("client exiting: %s", msg.Client.ID)
			g.removeClient(msg.Client)
			msg.Client.Close()
			break
		}
		go g.commander.HandleCommand(cmd, msg.Client)
	}
}

// handleClientPruning removes closed clients
func (g *Game) handleClientPruning() {
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

// broadcast sends a message to all clients:
func (g *Game) broadcast(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	for _, client := range g.clients {
		client.Write(msg)
	}
}

// removeClient removes a client:
func (g *Game) removeClient(client *Client) {
	g.log("before removeClient (client:%s) len(clients): %v", client.ID, len(g.clients))
	delete(g.clients, client.ID)
	g.log("after removeClient (client:%s) len(clients): %v", client.ID, len(g.clients))
}

// log is for logging a message
func (g *Game) log(message string, args ...interface{}) {
	msg := fmt.Sprintf("game:%s | %s\n", g.ID, message)
	log.Printf(msg, args...)
}
