package data

// Player is a player
type Player struct {
	Name string
}

// NewPlayer creates a player
func NewPlayer(name string) *Player {
	return &Player{Name: name}
}