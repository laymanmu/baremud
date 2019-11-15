package network

import (
	"fmt"
	"net"
	"strings"
	"baremud/data"
)

// Server is a network server
type Server struct {
	Clients  map[string]*Client
	Inbox    chan *Message
	Port     string
	Listener net.Listener
	World    *data.World
}

// NewServer creates a new network Server
func NewServer(port string, world *data.World) *Server {
	clients := make(map[string]*Client)
	inbox := make(chan *Message)
	return &Server{Clients: clients, Port: port, Inbox: inbox, World: world}
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
	message = fmt.Sprintf("%s\r\n", message)
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
		s.Inbox <- NewMessage(client.Addr, message)
		if message == "exit" {
			break
		}
	}
	fmt.Printf("DEBUG: finished handleConnection for %s\n", client.Player.Name)
}

// handleInbox will handle the network inbox messages
func (s *Server) handleInbox() {
	for {
		msg := <-s.Inbox
		fmt.Printf("network inbox | from: %v | to: %v | message: %v\n", msg.From, msg.To, msg.Message)
		client := s.Clients[msg.From]
		name := client.Player.Name
		if msg.Message == "exit" {
			s.removeClient(client)
			message := fmt.Sprintf("%s disconnected", name)
			s.Broadcast(message)
		} else if msg.Message == "look" {
			message := client.Player.Room.Look()
			s.Send(client.Addr, message)
		} else if len(msg.To) > 0 {
			message := fmt.Sprintf("%s says: %s\r\n", name, msg.Message)
			s.Send(msg.To, message)
		} else {
			message := fmt.Sprintf("%s says: %s\r\n", name, msg.Message)
			s.Broadcast(message)
		}
	}
}

func (s *Server) handleLogin(client *Client) error {
	client.Write("\r\nwhat is your name? ")
	data, err := client.Reader.ReadString('\n')
	if err != nil {
		return err
	}
	name := strings.TrimSpace(string(data))
	
	s.Send(client.Addr, fmt.Sprintf("\r\nWelcome, %s\r\n", name))
	client.Write("\r\nwhat is your password? ")
	data, err = client.Reader.ReadString('\n')
	if err != nil {
		return err
	}
	_ = strings.TrimSpace(string(data))
	// todo: auth

	client.Player.Name = name
	client.Player.EnterGate(s.World.GetStartGate())
	s.World.AddPlayer(client.Player)
	roomMessage := client.Player.Room.Look()
	client.Write(fmt.Sprintf("\r\n%s\r\n", roomMessage))
	fmt.Printf("DEBUG: finished handleLogin for %s\n", name)
	return nil
}

func (s *Server) removeClient(client *Client) {
	s.World.RemovePlayer(client.Player)
	client.Conn.Close()
	delete(s.Clients, client.Addr)
	fmt.Printf("removed network client connection: %s\n", client.Addr)
}