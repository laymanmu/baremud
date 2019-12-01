package main

import "sort"

// State is state
type State struct {
	ID        string
	TickCount int
	Players   map[string]*Player
	log       Logger
}

// NewState creates a state
func NewState() *State {
	id := NewID("state")
	log := NewLogger(id)
	players := make(map[string]*Player)
	return &State{id, 0, players, log}
}

// SortPlayersBy sorts players by a resource
func (s *State) SortPlayersBy(resourceName string) []*Player {
	defer (Track("sortPlayersBy", s.log))()
	names := make([]string, 0, len(s.Players))
	for name := range s.Players {
		names = append(names, name)
	}

	sort.Slice(names, func(x, y int) bool {
		l := s.Players[names[x]]
		r := s.Players[names[y]]
		return l.Resources[resourceName].Value > r.Resources[resourceName].Value
	})

	players := []*Player{}
	for _, name := range names {
		p := s.Players[name]
		players = append(players, p)
	}
	return players
}
