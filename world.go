package main

import (
	"fmt"
)

// World is a world
type World struct {
	StartRoom  *Room
	server     *Server
	rooms      map[string]*Room
	fromServer chan *Message
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

	fromServer := make(chan *Message)
	server := NewServer(":2323", fromServer, startRoom)

	return &World{StartRoom: startRoom, server: server, rooms: rooms, fromServer: fromServer}
}

// Start starts a world
func (w *World) Start() {
	w.server.Start()
	PrintMessageTypeValues()
	go w.handleServerMessages(w.fromServer)
}

// handleServerMessages handles messages from the server
func (w *World) handleServerMessages(messages <-chan *Message) {
	for {
		message := <-messages
		fmt.Printf("w | got msg: %+v\n", message)
		switch message.Type {
		case ClientLookMessage:
			message.Client.Write(message.Client.room.Look())
			fmt.Printf("w | handled msg: %+v\n", message)
		case ClientEnterMessage:
			message.Client.EnterGate(message.Message)
			fmt.Printf("w | handled msg: %+v\n", message)
		case ClientChatMessage:
			message.Client.room.Broadcast(message.Message)
			fmt.Printf("w | handled msg: %+v\n", message)
		case ClientStartedMessage:
			w.server.broadcast(fmt.Sprintf("%s joined", message.Client.Name))
			fmt.Printf("w | handled msg: %+v\n", message)
		case ClientStoppedMessage:
			w.server.broadcast(fmt.Sprintf("%s left", message.Client.Name))
			fmt.Printf("w | handled msg: %+v\n", message)
		default:
			fmt.Printf("w | unhandled msg | %+v\n", message)
		}
	}
}
