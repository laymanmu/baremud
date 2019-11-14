package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

// NetClient is a network client
type NetClient struct {
	Addr   string
	Conn   net.Conn
	Reader *bufio.Reader
	IsDone bool
}

// NewNetClient creates a network client
func NewNetClient(conn net.Conn) *NetClient {
	addr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	return &NetClient{Addr: addr, Conn: conn, Reader: reader, IsDone: false}
}

// Write will write a message to a client
func (c *NetClient) Write(message string) {
	c.Conn.Write([]byte(message))
}

// NetMessage is a network message
type NetMessage struct {
	From    string
	Message string
}

// NetServer is a network server
type NetServer struct {
	NetClients map[string]NetClient
	Port       string
	Listener   net.Listener
	Inbox      chan NetMessage
}

// NewNetServer creates a new NetServer
func NewNetServer(port string) *NetServer {
	clients := make(map[string]NetClient)
	inbox := make(chan NetMessage)
	return &NetServer{NetClients: clients, Port: port, Inbox: inbox}
}

// Start will start the server
func (s *NetServer) Start() {
	go s.listen()
	go s.handleInbox()
}

// DropClient will disconnect and remove a client
func (s *NetServer) DropClient(addr string) {
	client := s.NetClients[addr]
	fmt.Printf("dropping client: %s\n", client.Addr)
	client.Conn.Close()
	delete(s.NetClients, client.Addr)
	s.Inbox <- NetMessage{From: client.Addr, Message: "exit"}
}

// HandleInbox will handle the inbox messages
func (s *NetServer) handleInbox() {
	for {
		msg := <-s.Inbox
		fmt.Printf("inbox | from: %s | message: %s\n", msg.From, msg.Message)
		if msg.Message == "exit" {
			message := fmt.Sprintf("%s disconnected", msg.From)
			s.Broadcast(message)
		}
	}
}

// Broadcast will send a message to all clients
func (s *NetServer) Broadcast(message string) {
	fmt.Printf("broadcast: %s\n", message)
	for _, client := range s.NetClients {
		client.Write(message)
	}
}

// listen will loop and accept client connections.
func (s *NetServer) listen() {
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
			client := NewNetClient(conn)
			s.NetClients[client.Addr] = *client
			fmt.Printf("accepted client connection: %s\n", client.Addr)
			go s.handleConnection(client)
		}
	}
}

// handleConnection handles a new client connection
func (s *NetServer) handleConnection(client *NetClient) {
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
		s.Inbox <- NetMessage{From: client.Addr, Message: message}
	}
}
