package main

import (
	"fmt"
	"regexp"
	"strings"
)

// Resource is a resource
type Resource struct {
	ID    string
	Name  string
	Abbrv string
	Value int
	Max   int
	Min   int
	Delta int
}

// NewResource creates a resource
func NewResource(name, abbrv string, value, delta int) *Resource {
	return &Resource{NewID("resource"), name, abbrv, value, value, 0, delta}
}

// String returns a string
func (r *Resource) String() string {
	return fmt.Sprintf("%s:%v", r.Abbrv, r.Value)
}

// Player is a player
type Player struct {
	ID             string
	Name           string
	Resources      map[string]*Resource
	PromptTemplate string
	client         *Client
	log            Logger
}

// NewPlayer creates a player
func NewPlayer(name string, client *Client) *Player {
	id := NewID("player")
	log := NewLogger(id)
	resources := make(map[string]*Resource)
	resources["health"] = NewResource("health", "hp", 100, 5)
	resources["energy"] = NewResource("energy", "ep", 100, 20)
	resources["hunger"] = NewResource("hunger", "fp", 100, -1)
	resources["experience"] = NewResource("experience", "xp", 0, 0)
	resources["experience"].Max = 999999999
	promptTemplate := "[xp:%xp|fp:%fp|hp:%hp|ep:%ep]"
	return &Player{id, name, resources, promptTemplate, client, log}
}

// Update updates the player
func (p *Player) Update(game *Game) {
	p.updateResources()
}

// findResource returns a resource from an abbreviation
func (p *Player) findResource(abbrv string) *Resource {
	for _, resource := range p.Resources {
		if resource.Abbrv == abbrv {
			return resource
		}
	}
	return nil
}

// BuildPrompt returns a players prompt string
func (p *Player) BuildPrompt() string {
	rp := regexp.MustCompile("%[a-z]+")
	vars := Uniq(rp.FindAllString(p.PromptTemplate, -1))
	vals := make(map[string]string, len(vars))
	for _, v := range vars {
		abbrv := strings.Replace(v, "%", "", 1)
		res := p.findResource(abbrv)
		vals[v] = fmt.Sprintf("%v", res.Value)
	}
	prompt := p.PromptTemplate[:]
	for k, v := range vals {
		prompt = strings.Replace(prompt, k, v, -1)
	}
	return prompt
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
