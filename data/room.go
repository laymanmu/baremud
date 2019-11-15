package data

import (
	"bytes"
	"fmt"
)

// Room is a room
type Room struct {
	Name string
	Desc string
	Gates map[string]*Gate
	Players map[string]*Player
}

// NewRoom creates a room
func NewRoom(name, desc string) *Room {
	gates := make(map[string]*Gate, 2)
	players := make(map[string]*Player, 2)
	return &Room{Name: name, Desc: desc, Gates: gates, Players: players}
}

// AddGate adds a gate to a room
func (r *Room) AddGate(gate *Gate) {
	r.Gates[gate.Name] = gate
}
// RemoveGate removes a gate from the room
func (r *Room) RemoveGate(gate *Gate) {
	delete(r.Gates, gate.Name)
}

// AddPlayer adds a player to the room
func (r *Room) AddPlayer(player *Player) {
	r.Players[player.Name] = player
}
// RemovePlayer removes a player from the room
func (r *Room) RemovePlayer(player *Player) {
	delete(r.Players, player.Name)
}

// Look returns a string showing what is in this room
func (r *Room) Look() string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("-= %s =-\r\n", r.Name))
	buf.WriteString(fmt.Sprintf("%s\r\n", r.Desc))
	buf.WriteString(fmt.Sprintf("Gates: %v\r\n", r.Gates))
	buf.WriteString(fmt.Sprintf("Players: %v\r\n", r.Players))
	return buf.String()
}