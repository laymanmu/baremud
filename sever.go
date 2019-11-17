package main

import (
	"fmt"
	"net"

	"github.com/google/uuid"
)

// Server is a server
type Server struct {
	ID        string
	port      string
	clients   map[string]*Client
	messages  chan<- *Message
	listener  net.Listener
	startRoom *Room
}

// NewServer creates a server
func NewServer(port string, messages chan<- *Message, startRoom *Room) *Server {
	ID := uuid.New().String()
	clients := make(map[string]*Client)
	return &Server{ID, port, clients, messages, nil, startRoom}
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
		panic(err)
	}
	fmt.Printf("s | listening at %s\n", s.port)
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			fmt.Printf("s | client accept failed: %+v\n", err)
		} else {
			messages := make(chan *Message)
			client := NewClient(conn, messages, s.startRoom)
			s.clients[client.ID] = client
			fmt.Printf("s | added client: %s\n", client.ID)
			go s.handleClientMessages(client, messages)
		}
	}
}

// handleClientMessages handles messaging from one client
func (s *Server) handleClientMessages(client *Client, messages <-chan *Message) {
	for {
		message := <-messages
		fmt.Printf("s | got msg: %+v\n", message)
		switch message.Type {
		case ErrorMessage:
			message.Client.Write(message.Message)
			fmt.Printf("s | handled msg: %+v\n", message)
		case ClientStoppedMessage:
			delete(s.clients, client.ID)
			s.messages <- message
		default:
			s.messages <- message
		}
	}
}

// broadcast will send a message to all clients:
func (s *Server) broadcast(message string) {
	for _, client := range s.clients {
		client.Write(message)
	}
}
