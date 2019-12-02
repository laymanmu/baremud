package main

// State is state
type State struct {
	ID            string
	TickCount     int
	players       map[string]*Player
	places        map[string]*Place
	placePlayers  map[string][]*Player
	playerPlace   map[string]*Place
	clientPlayer  map[string]*Player
	log           Logger
	startingPlace *Place
}

// NewState creates a state
func NewState() *State {
	id := NewID("state")
	log := NewLogger(id)
	players := make(map[string]*Player)
	places := make(map[string]*Place)

	// id lookups:
	placePlayers := make(map[string][]*Player)
	playerPlace := make(map[string]*Place)
	clientPlayer := make(map[string]*Player)

	// build a state with a nil starting place:
	state := &State{id, 0, players, places, placePlayers, playerPlace, clientPlayer, log, nil}

	// add the starting place:
	state.startingPlace = NewPlace("Starting Place", "The starting place", state)
	places[state.startingPlace.ID] = state.startingPlace

	return state
}

// StartingPlace returns a starting place
func (s *State) StartingPlace() *Place {
	return s.startingPlace
}

// Add adds a player and/or place
func (s *State) Add(player *Player, place *Place) {
	defer (Track("Add", s.log))()
	s.players[player.ID] = player
	s.places[place.ID] = place
	s.clientPlayer[player.client.ID] = player
	s.MovePlayer(player, place)
}

// RemovePlayer removes a player
func (s *State) RemovePlayer(player *Player) {
	defer (Track("RemovePlayer", s.log))()
	s.RemovePlayerFromPlace(player, s.GetPlace(player))
	delete(s.clientPlayer, player.client.ID)
	delete(s.players, player.ID)
}

// RemovePlace removes a place
func (s *State) RemovePlace(place *Place) {
	defer (Track("RemovePlace", s.log))()
	delete(s.placePlayers, place.ID)
	delete(s.places, place.ID)
}

// GetPlace returns the place that the player is in
func (s *State) GetPlace(player *Player) *Place {
	defer (Track("GetPlace", s.log))()
	return s.playerPlace[player.ID]
}

// GetPlayers returns the players in a place
func (s *State) GetPlayers(place *Place) []*Player {
	defer (Track("GetPlayers", s.log))()
	return s.placePlayers[place.ID]
}

// GetAllPlayers returns all players
func (s *State) GetAllPlayers() []*Player {
	defer (Track("GetAllPlayers", s.log))()
	players := make([]*Player, len(s.players))
	i := 0
	for _, player := range s.players {
		players[i] = player
		i++
	}
	return players
}

// GetAllPlaces returns all places
func (s *State) GetAllPlaces() []*Place {
	n := "GetAllPlaces"
	defer (Track(n, s.log))()
	places := make([]*Place, len(s.places))
	for _, place := range s.places {
		places = append(places, place)
	}
	return places
}

// MovePlayer moves a player to a place
func (s *State) MovePlayer(player *Player, place *Place) {
	defer (Track("MovePlayer", s.log))()
	oldPlace := s.GetPlace(player)
	if oldPlace != nil {
		s.RemovePlayerFromPlace(player, oldPlace)
	}
	s.AddPlayerToPlace(player, place)
}

// RemovePlayerFromPlace removes a player from a place
func (s *State) RemovePlayerFromPlace(player *Player, place *Place) {
	defer (Track("RemovePlayerFromPlace", s.log))()
	delete(s.playerPlace, player.ID)
	index := s.indexOfPlacePlayers(place, player)
	if index == -1 {
		return
	}
	// copy last element into slot being removed, then remove last slot:
	players := s.placePlayers[place.ID]
	players[index] = players[len(players)-1]
	s.placePlayers[place.ID] = players[:len(players)-1]
}

// AddPlayerToPlace adds a player to a place
func (s *State) AddPlayerToPlace(player *Player, place *Place) {
	defer (Track("AddPlayerToPlace", s.log))()
	if place == nil {
		place = s.StartingPlace()
	}
	s.playerPlace[player.ID] = place
	index := s.indexOfPlacePlayers(place, player)
	if index != -1 {
		return
	}
	s.placePlayers[place.ID] = append(s.placePlayers[place.ID], player)
}

func (s *State) indexOfPlacePlayers(place *Place, player *Player) int {
	defer (Track("indexOfPlacePlayers", s.log))()
	// types:
	s.log("         player type: %T", player)
	s.log("          place type: %T", place)
	s.log("         place value: %v", place)
	s.log("       place.ID type: %T", place.ID)
	s.log("s.placePlayers  type: %T", s.placePlayers)
	// values:
	s.log("        player value: %v", player)
	s.log("         place value: %v", place)
	s.log("s.placePlayers value: %v", s.placePlayers)
	s.log("      place.ID value: %v", place.ID)
	players := s.placePlayers[place.ID]
	if players == nil {
		return -1
	}
	for i, p := range players {
		if p.ID == player.ID {
			return i
		}
	}
	return -1
}
