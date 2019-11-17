package main

import (
	"fmt"
)

// World is a world
type World struct {
	StartRoom  *Room
	server     *Server
	rooms      map[string]*Room
	fromServer chan interface{}
}

// NewWorld creates a world
func NewWorld() *World {
	rooms := make(map[string]*Room)

	cell := NewRoom("Cell", "A small cell")
	hall := NewRoom("Hall", "A dark hall")

	cell.Gates["door"] = hall
	hall.Gates["door"] = cell

	rooms[cell.ID] = cell
	rooms[hall.ID] = hall

	startRoom := cell

	fromServer := make(chan interface{})
	server := NewServer(":2323", fromServer, startRoom)

	return &World{StartRoom: startRoom, server: server, rooms: rooms, fromServer: fromServer}
}

// Start starts a world
func (w *World) Start() {
	w.server.Start()
	go w.handleServerMessages(w.fromServer)
}

// handleServerMessages handles messages from the server
func (w *World) handleServerMessages(messages <-chan interface{}) {
	for {
		m := <-messages
		switch message := m.(type) {
		case *ClientLookMessage:
			message.Client.Write(message.Client.room.Look())
		case *ClientEnterMessage:
			message.Client.EnterGate(message.Args[0])
		default:
			fmt.Printf("w | unhandled %T | %+v\n", message, message)
		}
	}
}
