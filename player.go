package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Player is a player
type Player struct {
	ID             string
	Name           string
	Resources      map[string]*Resource
	PromptTemplate string
	client         *Client
	log            Logger
	Place          *Place
}

// NewPlayer creates a player
func NewPlayer(name string, client *Client, place *Place) *Player {
	id := NewID("player")
	log := NewLogger(id)
	resources := make(map[string]*Resource)
	resources["health"] = NewResource("health", "hp", 100, 5)
	resources["energy"] = NewResource("energy", "ep", 100, 20)
	resources["hunger"] = NewResource("hunger", "fp", 100, -1)
	resources["experience"] = NewResource("experience", "xp", 0, 0)
	resources["experience"].Max = 999999999
	promptTemplate := "[xp:%xp|fp:%fp|hp:%hp|ep:%ep]"
	player := &Player{id, name, resources, promptTemplate, client, log, nil}
	place.AddPlayer(player)
	return player
}

// Update updates the player
func (p *Player) Update(game *Game) {
	for _, resource := range p.Resources {
		resource.Update()
	}
}

// BuildPrompt returns a players prompt string
func (p *Player) BuildPrompt() string {
	rp := regexp.MustCompile("%[a-z]+")
	vars := Uniq(rp.FindAllString(p.PromptTemplate, -1))
	vals := make(map[string]string, len(vars))
	for _, v := range vars {
		name := strings.Replace(v, "%", "", 1)
		if res, ok := p.Resources[name]; ok {
			vals[v] = fmt.Sprintf("%v", res.Value)
		}
	}
	prompt := p.PromptTemplate[:]
	for k, v := range vals {
		prompt = strings.Replace(prompt, k, v, -1)
	}
	return prompt
}

// MoveTo moves a player to a place
func (p *Player) MoveTo(place *Place) {
	if p.Place != nil {
		p.Place.RemovePlayer(p)
	}
	place.AddPlayer(p)
}

// Look sends a look response to the client
func (p *Player) Look() {
	p.client.Write(p.Place.Look(p))
}
