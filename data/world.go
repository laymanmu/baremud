package data

// World is a world
type World struct {
	Rooms map[string]*Room
	Players map[string]*Player
	Gates map[string]*Gate
}

// NewWorld creates a world
func NewWorld() *World {
	rooms := make(map[string]*Room, 10)
	players := make(map[string]*Player, 10)
	gates := make(map[string]*Gate, 10)
	w := &World{Rooms:rooms, Players:players, Gates:gates}

	cell := w.CreateRoom("Cell", "A small cell.")
	hall := w.CreateRoom("Hall", "A dark hall.")

	cell.AddGate(NewGate("Door", hall))
	hall.AddGate(NewGate("Door", cell))

	return w
}

// CreateRoom creates a room in the world
func (w *World) CreateRoom(name, desc string) *Room {
	room := NewRoom(name, desc)
	w.Rooms[room.Name] = room
	return room
}

// GetStartGate returns a gate to use go get to a starting room
func (w *World) GetStartGate() *Gate {
	for _, room := range w.Rooms {
		for _, gate := range room.Gates {
			return gate
		}
	}
	return nil
}

// RemovePlayer will remove a player from the world
func (w *World) RemovePlayer(player *Player) {
	player.Room.RemovePlayer(player)
	delete(w.Players, player.Name)
}

// AddPlayer will add a player to the world
func (w *World) AddPlayer(player *Player) {
	w.Players[player.Name] = player
}












