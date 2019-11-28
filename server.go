package main

import (
	"fmt"
	"log"
	"net"
)

// Server is a server
type Server struct {
	ID          string
	newClients  chan<- *Client
	clientInput chan<- *ClientInputMessage
	port        string
	listener    net.Listener
}

// NewServer creates a server
func NewServer(port string, newClients chan<- *Client, clientInput chan<- *ClientInputMessage) *Server {
	return &Server{NewID(), newClients, clientInput, port, nil}
}

// Start starts the server
func (s *Server) Start() {
	go s.listen()
}

// listen listens for connections
func (s *Server) listen() {
	var err error
	s.listener, err = net.Listen("tcp", s.port)
	if err != nil {
		s.log("listen() failed to start: %+v", err)
		panic(err)
	}
	s.log("listening at %s\n", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.log("client accept connection failed: %+v", err)
		} else {
			client := NewClient(conn, s.clientInput)
			s.newClients <- client
			s.log("accepted connection from client: %s", client.ID)
		}
	}
}

// log is for logging a message
func (s *Server) log(message string, args ...interface{}) {
	msg := fmt.Sprintf("server:%s | %s\n", s.ID, message)
	log.Printf(msg, args...)
}
