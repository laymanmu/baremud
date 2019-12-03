package main

import "strings"

// Place is a place
type Place struct {
	ID      string
	Name    string
	Desc    string
	Players []*Player
	log     Logger
}

// NewPlace creates a place
func NewPlace(name, desc string) *Place {
	id := NewID("place")
	log := NewLogger(id)
	players := []*Player{}
	return &Place{id, name, desc, players, log}
}

// Look looks at a place
func (p *Place) Look(looker *Player) string {
	playerNames := []string{}
	for _, player := range p.Players {
		if player.ID != looker.ID {
			playerNames = append(playerNames, player.Name)
		}
	}
	var sb strings.Builder
	sb.WriteString(p.Name)
	sb.WriteString("\r\n")
	sb.WriteString(p.Desc)
	sb.WriteString("\r\n")
	sb.WriteString("players: [")
	sb.WriteString(strings.Join(playerNames, ", "))
	sb.WriteString("]")
	sb.WriteString("\r\n")
	return sb.String()
}

// IsInRoom checks if a player is in this place
func (p *Place) IsInRoom(player *Player) bool {
	for _, roomPlayer := range p.Players {
		if roomPlayer.ID == player.ID {
			return true
		}
	}
	return false
}

// AddPlayer adds a player
func (p *Place) AddPlayer(player *Player) {
	player.Place = p
	if p.IsInRoom(player) {
		return
	}
	p.Players = append(p.Players, player)
}

// RemovePlayer removes a player
func (p *Place) RemovePlayer(player *Player) {
	player.Place = nil
	if !p.IsInRoom(player) {
		return
	}
	players := make([]*Player, len(p.Players)-1)
	for _, roomPlayer := range p.Players {
		if roomPlayer.ID != player.ID {
			players = append(players, roomPlayer)
		}
	}
	p.Players = players
}
