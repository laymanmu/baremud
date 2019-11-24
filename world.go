package main

import (
	"fmt"
	"strings"

	au "github.com/logrusorgru/aurora"
)

// World is a world
type World struct {
	StartRoom  *Room
	commander  *Commander
	server     *Server
	rooms      map[string]*Room
	fromServer chan interface{}
	commands   chan *CommandMessage
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

	fromServer := make(chan interface{})
	server := NewServer(":2323", fromServer)

	commander := NewCommander()
	commands := make(chan *CommandMessage)

	return &World{startRoom, commander, server, rooms, fromServer, commands}
}

// Start starts a world
func (w *World) Start() {
	s := "w | Start"
	defer Trace(s, "ended")
	Trace(s, "started")
	w.server.Start()
	go w.handleServerMessages()
	go w.handleCommandMessages()
}

func (w *World) handleCommandMessages() {
	s := "w | handleCommandMessages"
	defer Trace(s, "ended")
	Trace(s, "started")
	for {
		Trace(s, "polling")
		message := <-w.commands
		Trace(s, "got: %v", message)
		switch message.Command {
		case "look":
			w.commander.HandleLook(message)
		case "enter":
			w.commander.HandleEnter(message)
		case "exit":
			w.commander.HandleExit(message)
		case "say":
			w.commander.HandleSay(message)
		case "mkroom":
			w.commander.HandleMkroom(message)
		default:
			Trace(s, "unhandled command: %s", message.Command)
		}
	}
}

// handleServerMessages handles messages from the server
func (w *World) handleServerMessages() {
	for {
		message := <-w.fromServer
		switch msg := message.(type) {
		case *ErrorMessage:
			msg.Meta.Client.Write(msg.Message)
		case *ClientJoinMessage:
			player := msg.Meta.Client.Player
			player.Room = w.StartRoom
			player.Room.Clients[msg.Meta.Client.ID] = msg.Meta.Client
			msg.Meta.Client.Write(w.StartRoom.Look())
		case *ClientLeaveMessage:
			player := msg.Meta.Client.Player
			player.Room = w.StartRoom
			delete(player.Room.Clients, msg.Meta.Client.ID)
			player.Room.Broadcast("%s left", player.Name)
		case *InputMessage:
			fields := strings.Fields(msg.Input)
			command := fields[0]
			args := fields[1:]
			if w.commander.IsCommand(command) {
				w.commands <- NewCommandMessage(msg.Meta.Client, command, args)
			} else {
				msg.Meta.Client.Write("unknown command: %s. try: help", command)
			}
		default:
			Trace("w | handleServerMessages", "unhandled message type: %T, inspect: %+v", msg, msg)
		}
	}
}
