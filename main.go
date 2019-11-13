package main

import "fmt"

func main() {
	server := NewNetServer(":2323")
	server.Start()

	block := make(chan bool)
	<- block
	fmt.Printf("stopped\n")
}