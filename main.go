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
		stamp := time.Now().Format("2006/01/02 03:04:05.000")
		entry := fmt.Sprintf("%s %-15s | %s\n", stamp, source, message)
		fmt.Printf(entry, args...)
	}
}

// Track logs a start msg and returns a function to defer for the stop msg
// Example: defer (Track("myFunc", my.log))()
func Track(functionName string, log Logger) func() {
	id := NewID("track")
	log("%s | %s started", id, functionName)
	return func() {
		log("%s | %s stopped", id, functionName)
	}
}

// Uniq returns a unique list of strings:w
func Uniq(list []string) []string {
	keys := make(map[string]bool)
	uniq := []string{}
	for _, entry := range list {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			uniq = append(uniq, entry)
		}
	}
	return uniq
}

func main() {
	game := NewGame()
	game.Start()
}
