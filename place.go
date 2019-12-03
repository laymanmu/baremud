package main

import "strings"

// Place is a place
type Place struct {
	ID      string
	Name    string
	Desc    string
	players []*Player
	log     Logger
}

// NewPlace creates a place
func NewPlace(name, desc string) *Place {
	id := NewID("place")
	log := NewLogger(id)
	players := []*Player{}
	log("starting with players: %v", players)
	return &Place{id, name, desc, players, log}
}

// Look looks at a place
func (p *Place) Look(looker *Player) string {
	playerNames := []string{}
	for _, player := range p.players {
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
	for _, roomPlayer := range p.players {
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
	p.players = append(p.players, player)
	p.log("players after adding: %s", p.PlayerIDs())
}

// RemovePlayer removes a player
func (p *Place) RemovePlayer(player *Player) {
	player.Place = nil
	if !p.IsInRoom(player) {
		return
	}
	players := []*Player{}
	for _, roomPlayer := range p.players {
		if roomPlayer.ID != player.ID {
			players = append(players, roomPlayer)
		}
	}
	p.players = players
	p.log("players after removing: %s", p.PlayerIDs())
}

// PlayerIDs returns a list of Player IDs
func (p *Place) PlayerIDs() []string {
	ids := make([]string, len(p.players))
	for i, player := range p.players {
		ids[i] = player.ID
	}
	return ids
}
