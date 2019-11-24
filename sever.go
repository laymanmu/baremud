package main

import (
	"fmt"
	"net"
)

// Server is a server
type Server struct {
	ID       string
	port     string
	clients  map[string]*Client
	Upstream chan<- interface{}
	listener net.Listener
}

// NewServer creates a server
func NewServer(port string, upstream chan<- interface{}) *Server {
	clients := make(map[string]*Client)
	return &Server{GetID(), port, clients, upstream, nil}
}

// Start starts the server
func (s *Server) Start() {
	go s.listen()
}

// listen listens for connections
func (s *Server) listen() {
	src := "s | listen"
	var err error
	s.listener, err = net.Listen("tcp", s.port)
	if err != nil {
		Trace(src, "listen failed: %+v", err)
		panic(err)
	}
	Trace(src, "listening at %s", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			Trace(src, "client accept failed: %+v", err)
		} else {
			downstream := make(chan interface{})
			client := NewClient(conn, downstream)
			s.clients[client.ID] = client
			go s.handleClientMessages(client, downstream)
			Trace(src, "added client: %s\n", client.ID)
		}
	}
}

// handleClientMessages handles messaging from one client
func (s *Server) handleClientMessages(client *Client, messages <-chan interface{}) {
	src := fmt.Sprintf("s | handleClientMessages | client:%s", client.ID)
	for {
		message := <-messages
		switch msg := message.(type) {
		case ErrorMessage:
			Trace(src, "ErrorMessage | id:%s | %s", msg.Meta.ID, msg.Message)
			msg.Meta.Client.Write(msg.Message)
		case ClientLeaveMessage:
			Trace(src, "ClientLeaveMessage | client: %s", msg.Meta.Client.ID)
			delete(s.clients, msg.Meta.Client.ID)
			s.Upstream <- msg
		default:
			s.Upstream <- msg
		}
	}
}

// broadcast will send a message to all clients:
func (s *Server) broadcast(message string) {
	Trace("s | broadcast | %s", message)
	for _, client := range s.clients {
		client.Write(message)
	}
}
