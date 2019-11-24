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
func (r *Room) Look() string {
	var gates []string
	var clients []string
	for name := range r.Gates {
		gates = append(gates, fmt.Sprintf("%s", au.Blue(name)))
	}
	for _, c := range r.Clients {
		clients = append(clients, fmt.Sprintf("%s", au.Green(c.Player.Name)))
	}

	name := fmt.Sprintf("%s %s %s", au.BrightBlack("-=["), au.BrightYellow(r.Name), au.BrightBlack("]=-\r\n"))
	desc := fmt.Sprintf("%s\r\n", au.White(r.Desc))
	gateNames := fmt.Sprintf("%s %s%s\r\n", au.BrightBlack("[gates:"), strings.Join(gates, ", "), au.BrightBlack("]"))
	clientNames := fmt.Sprintf("%s %s%s", au.BrightBlack("[clients:"), strings.Join(clients, ", "), au.BrightBlack("]"))

	var buf bytes.Buffer
	buf.WriteString("\r\n")
	buf.WriteString(name)
	buf.WriteString(desc)
	buf.WriteString(gateNames)
	buf.WriteString(clientNames)
	return buf.String()
}

// Broadcast a message to the room
func (r *Room) Broadcast(message string, m ...interface{}) {
	msg := fmt.Sprintf(message, m...)
	for _, client := range r.Clients {
		client.Write(msg)
	}
}
