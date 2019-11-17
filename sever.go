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
	messages  chan<- interface{}
	listener  net.Listener
	startRoom *Room
}

// NewServer creates a server
func NewServer(port string, messages chan<- interface{}, startRoom *Room) *Server {
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
			messages := make(chan interface{})
			client := NewClient(conn, messages, s.startRoom)
			s.clients[client.ID] = client
			fmt.Printf("s | added client: %+v\n", client)
			go s.handleConnection(client, messages)
		}
	}
}

// handleConnection handles messaging from one client
func (s *Server) handleConnection(client *Client, clientMessages <-chan interface{}) {
	for {
		m := <-clientMessages
		switch message := m.(type) {
		case *ErrorMessage:
			fmt.Printf("s | client %s error: %+v\n", client.ID, message)
			message.Client.Write(message.Message)
		case *ClientChatMessage:
			msg := fmt.Sprintf("%s says: %s", message.Client.Name, message.Message)
			s.broadcast(msg, nil)
		case *ClientLoggedOnMessage:
			msg := fmt.Sprintf("%s has joined", client.Name)
			s.broadcast(msg, nil)
			s.messages <- message
		case *ClientClosedMessage:
			delete(s.clients, client.ID)
			msg := fmt.Sprintf("%s has left", client.Name)
			s.broadcast(msg, nil)
			s.messages <- message

		case *ClientLookMessage:
			s.messages <- message
		case *ClientEnterMessage:
			s.messages <- message

		default:
			fmt.Printf("s | client %s sent unhandled %T msg: %+v\n", client.ID, message, message)
		}
	}
}

// broadcast will write a message to a list of clients
func (s *Server) broadcast(message string, clients map[string]*Client) {
	fmt.Printf("s | broadcasting: %s\n", message)
	if clients == nil {
		clients = s.clients
	}
	for _, client := range clients {
		client.Write(message)
	}
}
