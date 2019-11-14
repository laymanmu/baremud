package network

import (
	"fmt"
	"net"
	"strings"
)

// Server is a network server
type Server struct {
	Clients  map[string]*Client
	Inbox    chan *Message
	Port     string
	Listener net.Listener
}

// NewServer creates a new network Server
func NewServer(port string) *Server {
	clients := make(map[string]*Client)
	inbox := make(chan *Message)
	return &Server{Clients: clients, Port: port, Inbox: inbox}
}

// Start will start the network server
func (s *Server) Start() {
	go s.listen()
	go s.handleInbox()
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
	message = fmt.Sprintf("%s\n", message)
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
	fmt.Printf("listening for network connections on port %s\n", s.Port)
	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			fmt.Printf("failed to accept network client connection: %s\n", err)
		} else {
			client := NewClient(conn)
			s.Clients[client.Addr] = client
			fmt.Printf("added network client connection: %s\n", client.Addr)
			go s.handleConnection(client)
		}
	}
}

// handleConnection handles a new network client connection
func (s *Server) handleConnection(client *Client) {
	err := s.handleLogin(client)
	if err != nil {
		s.removeClient(client)
		return
	}
	for {
		data, err := client.Reader.ReadString('\n')
		if err != nil {
			s.Inbox <- NewMessage(client.Addr, "exit")
			break
		}
		message := strings.TrimSpace(string(data))
		if message == "exit" {
			break
		} else {
			s.Inbox <- NewMessage(client.Addr, message)
		}
	}
	s.removeClient(client)
}

// handleInbox will handle the network inbox messages
func (s *Server) handleInbox() {
	for {
		msg := <-s.Inbox
		fmt.Printf("network inbox | from: %s | message: %s\n", msg.From, msg.Message)
		if msg.Message == "exit" {
			message := fmt.Sprintf("%s disconnected", msg.From)
			s.Broadcast(message)
		} else {
			name := s.Clients[msg.From].Player.Name
			message := fmt.Sprintf("%s says: %s", name, msg.Message)
			s.Broadcast(message)
		}
	}
}

func (s *Server) handleLogin(client *Client) error {
	client.Write("what is your name? ")
	data, err := client.Reader.ReadString('\n')
	if err != nil {
		return err
	}
	name := strings.TrimSpace(string(data))
	
	client.Write("what is your password? ")
	data, err = client.Reader.ReadString('\n')
	if err != nil {
		return err
	}
	password := strings.TrimSpace(string(data))
	fmt.Printf("TODO: authenticate name: %s, password: %s\n", name, password)
	client.Player.Name = name
	message := fmt.Sprintf("Welcome, %s", client.Player.Name)
	s.Send(client.Addr, message)
	return nil
}

func (s *Server) removeClient(client *Client) {
	client.Conn.Close()
	delete(s.Clients, client.Addr)
	fmt.Printf("removed network client connection: %s\n", client.Addr)
	s.Inbox <- NewMessage(client.Addr, "exit")
}