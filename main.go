package main

func main() {
	world := NewWorld()
	world.Start()

	block := make(chan bool)
	<-block
}
