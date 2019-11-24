package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Trace will trace log a message
func Trace(source, msgFormat string, m ...interface{}) {
	stamp := time.Now().Format("15:04:05:06")
	prefix := fmt.Sprintf("%s | trace | %s", stamp, source)
	message := fmt.Sprintf(msgFormat, m...)
	fmt.Printf("%s | %s\n", prefix, message)
}

// GetID creates a unique id
func GetID() string {
	id := uuid.New().String()
	i := strings.IndexByte(id, '-')
	if i == -1 {
		return id
	}
	return id[:i]
}

var world *World

func main() {
	world = NewWorld()
	world.Start()
	block := make(chan bool)
	<-block
}
