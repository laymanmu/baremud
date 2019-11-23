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
	commands   chan *Message
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
	commands := make(chan *Message)

	return &World{startRoom, commander, server, rooms, fromServer, commands}
}

// Start starts a world
func (w *World) Start() {
	w.server.Start()
	go w.handleServerMessages()
	go w.handleCommandMessages()
}

// handleCommandMessages handles command messages from client input
func (w *World) handleCommandMessages() {
	for {

	}
}

// handleServerMessages handles messages from the server
func (w *World) handleServerMessages() {
	for {
		message := <-w.fromServer
		status := "handled"
		switch message.Type {
		case HelpMessage:
			message.Client.Write(w.commander.Help(message.Message))
		case ErrorMessage:
			message.Client.Write(message.Message)
		case InputMessage:
			command := message.Args[0]
			if w.commander.IsCommand(command) {
				w.commands <- NewMessage(CommandMessage, message.Client, command, message.Args[1:])
			} else {
				msg := fmt.Sprintf("unknonw command: %s", command)
				w.fromServer <- NewMessage(ErrorMessage, message.Client, msg, message.Args)
			}
		default:
			status = "unhandled"
		}
		fmt.Printf("w | %s %s\n", status, message.String())
	}
}
