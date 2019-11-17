package main

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/uuid"
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
	id := uuid.New().String()
	gates := make(map[string]*Room)
	clients := make(map[string]*Client)
	return &Room{ID: id, Name: name, Desc: desc, Gates: gates, Clients: clients}
}

// Look describes a room
func (r *Room) Look() string {
	var gates []string
	var clients []string
	for name := range r.Gates {
		gates = append(gates, name)
	}
	for _, c := range r.Clients {
		clients = append(clients, c.Name)
	}
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("-= %s =-\r\n", r.Name))
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
