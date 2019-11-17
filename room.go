package main

import (
	"bytes"
	"fmt"

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
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("-= %s =-\r\n", r.Name))
	buf.WriteString(fmt.Sprintf("%s\r\n", r.Desc))
	buf.WriteString(fmt.Sprintf("Gates:   %v\r\n", r.Gates))
	buf.WriteString(fmt.Sprintf("Clients: %v\r\n", r.Clients))
	return buf.String()
}

// Broadcast a message to the room
func (r *Room) Broadcast(message string) {
	for _, client := range r.Clients {
		client.Write(message)
	}
}
