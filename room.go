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
		clients = append(clients, fmt.Sprintf("%s", au.Green(c.Name)))
	}
	var buf bytes.Buffer

	buf.WriteString(fmt.Sprintf("%s %s %s\r\n", au.BrightBlack("-="), au.BrightYellow(r.Name), au.BrightBlack("=-")))
	buf.WriteString(fmt.Sprintf("%s\r\n", r.Desc))
	buf.WriteString(fmt.Sprintf("Gates:   %s\r\n", strings.Join(gates, ", ")))
	buf.WriteString(fmt.Sprintf("Clients: %s\r\n", strings.Join(clients, ", ")))
	return buf.String()
}

// Broadcast a message to the room
func (r *Room) Broadcast(message string) {
	for _, client := range r.Clients {
		client.Write(message)
	}
}
