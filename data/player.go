package data

// Player is a player
type Player struct {
	Name string
	Room *Room
}

// NewPlayer creates a player
func NewPlayer(name string) *Player {
	return &Player{Name: name}
}

// EnterGate moves the player through a gate to another room
func (p *Player) EnterGate(gate *Gate) {
	if p.Room != nil {
		p.Room.RemovePlayer(p)
	}
	p.Room = gate.To
	p.Room.Players[p.Name] = p
}
