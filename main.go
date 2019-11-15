package main

import (
	"baremud/network"
	"baremud/data"
)

func main() {
	world := data.NewWorld()
	server := network.NewServer(":2323", world)
	server.Start()

	block := make(chan bool)
	<-block
}
