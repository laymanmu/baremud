package main

import "fmt"

// Place is a place
type Place struct {
	ID    string
	Name  string
	Desc  string
	log   Logger
	state *State
}

// NewPlace creates a place
func NewPlace(name, desc string, state *State) *Place {
	id := NewID("place")
	log := NewLogger(id)
	place := &Place{id, name, desc, log, state}
	//
	state.places[id] = place
	state.placePlayers[id] = []*Player{}
	return place
}

// LookAt returns a look at this place
func (p *Place) LookAt(player *Player) string {
	f := "-= %s =-\r\n%s\r\n\r\nmore stuff...\r\n"
	return fmt.Sprintf(f, p.Name, p.Desc)
}

// Players returns the players in this room
func (p *Place) Players() []*Player {
	if players, ok := p.state.placePlayers[p.ID]; ok {
		return players
	}
	empty := []*Player{}
	p.state.placePlayers[p.ID] = empty
	return empty
}

// HasPlayer checks if a player is in this room
func (p *Place) HasPlayer(player *Player) bool {
	for _, player2 := range p.Players() {
		if player.ID == player2.ID {
			return true
		}
	}
	return false
}
