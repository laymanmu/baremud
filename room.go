package main

import (
	"bytes"
	"fmt"
	"strings"

	au "github.com/logrusorgru/aurora"
)

// Room is a room
type Room struct {
	ID      string
	Name    string
	Desc    string
	Gates   map[string]*Room
	Clients map[string]*Client
}

// NewRoom creates a room
func NewRoom(name, desc string) *Room {
	gates := make(map[string]*Room)
	clients := make(map[string]*Client)
	return &Room{GetID(), name, desc, gates, clients}
}

// Look describes a room
func (r *Room) Look(prompt string) string {
	var gates []string
	var clients []string
	for name := range r.Gates {
		gates = append(gates, fmt.Sprintf("%s", au.Blue(name)))
	}
	for _, c := range r.Clients {
		clients = append(clients, fmt.Sprintf("%s", au.Green(c.Name)))
	}

	name := fmt.Sprintf("%s %s %s", au.BrightBlack("-=["), au.BrightYellow(r.Name), au.BrightBlack("]=-"))
	title := fmt.Sprintf("%s %s\r\n", name, prompt)
	desc := fmt.Sprintf("%s\r\n", au.White(r.Desc))
	gateNames := fmt.Sprintf("%s %s%s\r\n", au.BrightBlack("[gates:"), strings.Join(gates, ", "), au.BrightBlack("]"))
	clientNames := fmt.Sprintf("%s %s%s", au.BrightBlack("[clients:"), strings.Join(clients, ", "), au.BrightBlack("]"))

	var buf bytes.Buffer
	buf.WriteString("\r\n")
	buf.WriteString(title)
	buf.WriteString(desc)
	buf.WriteString(gateNames)
	buf.WriteString(clientNames)
	return buf.String()
}

// Broadcast a message to the room
func (r *Room) Broadcast(message string) {
	for _, client := range r.Clients {
		client.Write(message)
	}
}
