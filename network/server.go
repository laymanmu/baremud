package network

import (
	"fmt"
	"net"
	"strings"
)

// Server is a network server
type Server struct {
	Clients  map[string]Client
	Port     string
	Listener net.Listener
	Inbox    chan Message
}

// NewServer creates a new NetServer
func NewServer(port string) *Server {
	clients := make(map[string]Client)
	inbox := make(chan Message)
	return &Server{Clients: clients, Port: port, Inbox: inbox}
}

// Start will start the server
func (s *Server) Start() {
	go s.listen()
	go s.handleInbox()
}

// DropClient will disconnect and remove a client
func (s *Server) DropClient(addr string) {
	client := s.Clients[addr]
	fmt.Printf("dropping client: %s\n", client.Addr)
	client.Conn.Close()
	delete(s.Clients, client.Addr)
	s.Inbox <- Message{From: client.Addr, Message: "exit"}
}

// Broadcast will send a message to all network clients
func (s *Server) Broadcast(message string) {
	fmt.Printf("broadcasting: %s\n", message)
	for addr := range s.Clients {
		s.Send(addr, message)
	}
}

// Send will send a message to a given network client
func (s *Server) Send(addr, message string) {
	client := s.Clients[addr]
	client.Write(message)
}

// listen will loop and accept client connections.
func (s *Server) listen() {
	var err error
	s.Listener, err = net.Listen("tcp", s.Port)
	if err != nil {
		panic(err)
	}
	fmt.Printf("listening for connections on port %s\n", s.Port)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept client connection: %s\n", err)
		} else {
			client := NewClient(conn)
			s.Clients[client.Addr] = *client
			fmt.Printf("accepted client connection: %s\n", client.Addr)
			go s.handleConnection(client)
		}
	}
}

// handleConnection handles a new network client connection
func (s *Server) handleConnection(client *Client) {
	defer client.Conn.Close()
	for {
		data, err := client.Reader.ReadString('\n')
		if err != nil {
			s.DropClient(client.Addr)
			break
		}
		message := strings.TrimSpace(string(data))
		if message == "exit" {
			s.DropClient(client.Addr)
			break
		}
		s.Inbox <- Message{From: client.Addr, Message: message}
	}
}

// handleInbox will handle the inbox messages
func (s *Server) handleInbox() {
	for {
		msg := <-s.Inbox
		fmt.Printf("inbox | from: %s | message: %s\n", msg.From, msg.Message)
		if msg.Message == "exit" {
			message := fmt.Sprintf("%s disconnected", msg.From)
			s.Broadcast(message)
		}
	}
}
