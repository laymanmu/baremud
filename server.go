package main

import (
	"net"
)

// Server is a server
type Server struct {
	ID          string
	newClients  chan<- *Client
	clientInput chan<- *ClientInputMessage
	port        string
	listener    net.Listener
	log         Logger
}

// NewServer creates a server
func NewServer(port string, newClients chan<- *Client, clientInput chan<- *ClientInputMessage) *Server {
	id := NewID("server")
	log := NewLogger(id)
	return &Server{id, newClients, clientInput, port, nil, log}
}

// Start starts the server
func (s *Server) Start() {
	defer (Track("Start", s.log))()
	go s.listen()
}

// listen listens for connections
func (s *Server) listen() {
	defer (Track("listen", s.log))()
	var err error
	s.listener, err = net.Listen("tcp", s.port)
	if err != nil {
		s.log("listen() failed to start per error: %+v", err)
		panic(err)
	}
	s.log("listening at %s", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.log("failed to accept client per error: %+v", err)
		} else {
			client := NewClient(conn, s.clientInput)
			s.log("accepted %s", client.ID)
			s.newClients <- client
		}
	}
}
