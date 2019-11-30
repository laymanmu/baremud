package main

import (
	"fmt"
	"strings"
)

// Resource is a resource
type Resource struct {
	ID    string
	Name  string
	Value int
	Max   int
	Min   int
	Delta int
}

// NewResource creates a resource
func NewResource(name string, value, delta int) *Resource {
	return &Resource{NewID("resource"), name, value, value, 0, 10}
}

// String returns a string
func (r *Resource) String() string {
	return fmt.Sprintf("[%s:%v/%v]", r.Name, r.Value, r.Max)
}

// Player is a player
type Player struct {
	ID        string
	Name      string
	Resources map[string]*Resource
	client    *Client
	log       Logger
}

// NewPlayer creates a player
func NewPlayer(name string, client *Client) *Player {
	id := NewID("player")
	log := NewLogger(id)
	resources := make(map[string]*Resource)
	return &Player{id, name, resources, client, log}
}

// Update updates the player
func (p *Player) Update(game *Game) {
	p.updateResources()
	p.client.Write(p.Stats())
	p.log(p.Stats())
}

// Stats returns a players stats
func (p *Player) Stats() string {
	b := strings.Builder{}
	b.WriteString(p.Name)
	b.WriteString(" | ")
	for _, res := range p.Resources {
		b.WriteString(res.String())
	}
	return b.String()
}

// updateResources applies each resource delta
func (p *Player) updateResources() {
	for _, res := range p.Resources {
		v := res.Value + res.Delta
		if v > res.Max {
			v = res.Max
		} else if v < res.Min {
			v = res.Min
		}
		res.Value = v
	}
}
