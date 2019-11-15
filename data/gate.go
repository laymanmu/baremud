package data

import (
	uuid "github.com/satori/go.uuid"
)

// Gate is a gateway to a room
type Gate struct {
	ID   string
	Name string
	To   *Room
}

// NewGate creates a gate
func NewGate(name string, to *Room) *Gate {
	id, _ := uuid.NewV4()
	return &Gate{Name: name, To: to, ID: id.String()}
}