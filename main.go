package main

import (
	"strings"

	"github.com/google/uuid"
)

// GetID creates a unique id
func GetID() string {
	id := uuid.New().String()
	i := strings.IndexByte(id, '-')
	if i == -1 {
		return id
	}
	return id[:i]
}

func main() {
	world := NewWorld()
	world.Start()

	block := make(chan bool)
	<-block
}
