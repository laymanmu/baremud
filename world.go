package main

import (
	"fmt"

	au "github.com/logrusorgru/aurora"
)

// World is a world
type World struct {
	StartRoom  *Room
	commander  *Commander
	server     *Server
	rooms      map[string]*Room
	fromServer chan *Message
}

// NewWorld creates a world
func NewWorld() *World {
	rooms := make(map[string]*Room)

	name := fmt.Sprintf("%s Room", au.Red("Red"))
	desc := fmt.Sprintf("A small, %s room", au.Red("red"))
	room := NewRoom(name, desc)

	name = fmt.Sprintf("%s Hall", au.Green("Green"))
	desc = fmt.Sprintf("A dark, %s hall", au.Green("green"))
	hall := NewRoom(name, desc)

	room.Gates["door"] = hall
	hall.Gates["door"] = room

	rooms[room.ID] = room
	rooms[hall.ID] = hall

	startRoom := room

	fromServer := make(chan *Message)
	server := NewServer(":2323", fromServer, startRoom)

	commander := NewCommander()

	return &World{startRoom, commander, server, rooms, fromServer}
}

// Start starts a world
func (w *World) Start() {
	w.server.Start()
	go w.handleServerMessages(w.fromServer)
}

// handleServerMessages handles messages from the server
func (w *World) handleServerMessages(messages <-chan *Message) {
	chatPrefix := fmt.Sprintf("%s", au.BrightBlue("[chat]"))
	serverPrefix := fmt.Sprintf("%s", au.BrightRed("[server]"))
	for {
		message := <-messages
		status := "handled"
		switch message.Type {
		case HelpMessage:
			message.Client.Write(w.commander.Help(message.Message))
		case ClientLookMessage:
			message.Client.Write(message.Client.room.Look(message.Client.Prompt()))
		case ClientEnterMessage:
			message.Client.EnterGate(message.Message)
		case ClientChatMessage:
			msg := fmt.Sprintf("%s %s: %s", chatPrefix, au.BrightGreen(message.Client.Name), au.Cyan(message.Message))
			message.Client.room.Broadcast(msg)
		case ClientStartedMessage:
			w.server.broadcast(fmt.Sprintf("%s %s joined", serverPrefix, message.Client.Name))
		case ClientStoppedMessage:
			w.server.broadcast(fmt.Sprintf("%s %s left", serverPrefix, message.Client.Name))
		case ClientMakeRoomMessage:
			name := message.Args[0]
			desc := message.Args[1]
			gate := message.Args[2]
			room := NewRoom(name, desc)
			message.Client.room.Gates[gate] = room
			room.Gates[gate] = message.Client.room
			message.Client.Write(message.Client.room.Look(message.Client.Prompt()))
		default:
			status = "unhandled"
		}
		fmt.Printf("w | %s %s\n", status, message.String())
	}
}
