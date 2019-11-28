package main

import (
	"fmt"
	"strings"
	"time"

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

// Log logs a message from a given source like: "client:abc"
func Log(source, message string) {
	stamp := time.Now().Format("2006/01/02 03:04:05")
	entry := fmt.Sprintf("%s %-15s | %s", stamp, source, message)
	fmt.Println(entry)
}

func main() {
	game := NewGame()
	game.Start()
}
