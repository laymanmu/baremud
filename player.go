package main

import (
	"fmt"

	au "github.com/logrusorgru/aurora"
)

// Player is a player
type Player struct {
	client *Client
	Name   string
	Room   *Room
}

// NewPlayer creates a player
func NewPlayer(client *Client, name string) *Player {
	return &Player{client, name, world.StartRoom}
}

// Prompt returns the players prompt string
func (p *Player) Prompt() string {
	left := au.BrightBlack("-=[")
	right := au.BrightBlack("]=-")
	alert := au.Green("ok")
	return fmt.Sprintf("%s%s%s", left, alert, right)
}

// EnterGate will enter a gate if it exists
func (p *Player) EnterGate(name string) {
	if newRoom, ok := p.Room.Gates[name]; ok {
		delete(p.Room.Clients, p.client.ID)
		for _, client := range p.Room.Clients {
			client.Write(fmt.Sprintf("%s left the room", au.Green(p.Name)))
		}
		for _, client := range newRoom.Clients {
			client.Write(fmt.Sprintf("%s entered the room", au.Green(p.Name)))
		}
		p.Room = newRoom
		p.Room.Clients[p.client.ID] = p.client
		p.client.Write(p.Room.Look())
	} else {
		p.client.Write(fmt.Sprintf("gate %s not found", name))
	}
}
