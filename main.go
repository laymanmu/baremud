package main

import (
	"strings"

	"github.com/google/uuid"
)

// NewID creates a unique id
func NewID() string {
	id := uuid.New().String()
	i := strings.IndexByte(id, '-')
	if i == -1 {
		return id
	}
	return id[:i]
}

func main() {
	game := NewGame()
	game.Start()
}
