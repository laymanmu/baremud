package main

import (
	"baremud/network"
)

func main() {
	server := network.NewServer(":2323")
	server.Start()
	block := make(chan bool)
	<-block
}
