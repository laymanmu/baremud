package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Logger is a log function
type Logger func(string, ...interface{})

// NewID creates a unique id like "client:1a6bb7"
func NewID(typeName string) string {
	id := uuid.New().String()
	i := strings.IndexByte(id, '-')
	return fmt.Sprintf("%s:%s", typeName, id[:i])
}

// NewLogger creates a log() function that includes the source
func NewLogger(source string) Logger {
	return func(message string, args ...interface{}) {
		stamp := time.Now().Format("2006/01/02 03:04:05")
		entry := fmt.Sprintf("%s %-15s | %s\n", stamp, source, message)
		fmt.Printf(entry, args...)
	}
}

// NewTracker logs started/stopped messages for a function call
func NewTracker(functionName string, log Logger) func() {
	return func() {
		log("%s started", functionName)
		defer log("%s stopped", functionName)
	}
}

func main() {
	game := NewGame()
	game.Start()
}
