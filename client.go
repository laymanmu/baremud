package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/google/uuid"
)

// Client is a client
type Client struct {
	ID       string
	Name     string
	room     *Room
	conn     net.Conn
	reader   *bufio.Reader
	messages chan<- interface{}
}

// NewClient creates a client
func NewClient(conn net.Conn, messages chan<- interface{}, room *Room) *Client {
	ID := uuid.New().String()
	Name := ""
	reader := bufio.NewReader(conn)
	client := &Client{ID, Name, room, conn, reader, messages}
	client.room.Clients[client.ID] = client
	go client.handleConnection()
	return client
}

// EnterGate will enter a gate
func (c *Client) EnterGate(name string) {
	if newRoom, ok := c.room.Gates[name]; ok {
		delete(c.room.Clients, c.ID)
		for _, client := range c.room.Clients {
			client.Write(fmt.Sprintf("%s left the room", c.Name))
		}
		for _, client := range newRoom.Clients {
			client.Write(fmt.Sprintf("%s entered the room", c.Name))
		}
		c.room = newRoom
		c.room.Clients[c.ID] = c
		c.Write(c.room.Look())
	}
}

// Write writes a message to the client
func (c *Client) Write(message string) {
	c.conn.Write([]byte(fmt.Sprintf("%s\r\n", message)))
}

// handleLogin handles the client login
func (c *Client) handleLogin() error {
	c.Write("what is your name?")
	data, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	c.Name = strings.TrimSpace(string(data))
	c.Write(fmt.Sprintf("Welcome, %s", c.Name))
	c.messages <- NewClientLoggedOnMessage(c)
	return nil
}

func (c *Client) closeConnection() {
	c.conn.Close()
	c.messages <- NewClientClosedMessage(c)
}

// handleConnection handles a new network client connection
func (c *Client) handleConnection() {
	defer c.closeConnection()

	err := c.handleLogin()
	if err != nil {
		return
	}

	c.Write(c.room.Look())

	for {
		data, err := c.reader.ReadString('\n')
		if err != nil {
			return
		}

		message := strings.TrimSpace(string(data))
		fields := strings.Fields(message)
		if len(fields) < 1 {
			continue
		}
		command := fields[0]
		args := fields[1:]

		switch command {
		case "exit":
			return
		case "look":
			c.messages <- NewClientLookMessage(c, args)
		case "enter":
			c.messages <- NewClientEnterMessage(c, args)
		case "say":
			msg := strings.Join(args, " ")
			c.messages <- NewClientChatMessage(c, msg)
		default:
			msg := fmt.Sprintf("unknown command: %s, args: %v", command, args)
			c.messages <- NewErrorMessage(c, msg)
		}
	}
}
