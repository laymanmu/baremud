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
	Writer *bufio.Writer
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

// HandleInbox will handle the inbox messages
func (s *NetServer) handleInbox() {
	for {
		msg := <-s.Inbox
		fmt.Printf("inbox | from: %s | message: %s\n", msg.From, msg.Message)
		if msg.Message == "exit" {
			fmt.Printf("removing client: %s\n", msg.From)
			delete(s.NetClients, msg.From)
			message := fmt.Sprintf("%s disconnected", msg.From)
			s.Broadcast(message)
		}
	}
}

// Broadcast will send a message to all clients
func (s *NetServer) Broadcast(message string) {
	fmt.Printf("broadcast: %s\n", message)
	for _, client := range s.NetClients {
		client.Writer.Write([]byte(message))
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
			addr := conn.RemoteAddr().String()
			reader := bufio.NewReader(conn)
			writer := bufio.NewWriter(conn)
			client := NetClient{Addr: addr, Conn: conn, Reader: reader, Writer: writer}
			s.NetClients[addr] = client
			fmt.Printf("accepted client connection: %s\n", addr)
			go s.HandleConnection(&client)
		}
	}
}

// HandleConnection handles a new client connection
func (s *NetServer) HandleConnection(client *NetClient) {
	for {
		data, err := client.Reader.ReadString('\n')
		if err != nil {
			fmt.Printf("client disconnected: %s\n", client.Addr)
			s.Inbox <- NetMessage{From: client.Addr, Message: "exit"}
		}
		message := strings.TrimSpace(string(data))
		s.Inbox <- NetMessage{From: client.Addr, Message: message}
	}
}
