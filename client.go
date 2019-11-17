package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	au "github.com/logrusorgru/aurora"
)

// Client is a client
type Client struct {
	ID       string
	Name     string
	room     *Room
	conn     net.Conn
	reader   *bufio.Reader
	messages chan<- *Message
}

// NewClient creates a client
func NewClient(conn net.Conn, messages chan<- *Message, room *Room) *Client {
	Name := ""
	reader := bufio.NewReader(conn)
	client := &Client{GetID(), Name, room, conn, reader, messages}
	client.room.Clients[client.ID] = client
	go client.handleConnection()
	return client
}

// EnterGate will enter a gate if it exists
func (c *Client) EnterGate(name string) {
	if newRoom, ok := c.room.Gates[name]; ok {
		delete(c.room.Clients, c.ID)
		for _, client := range c.room.Clients {
			client.Write(fmt.Sprintf("%s left the room", au.Green(c.Name)))
		}
		for _, client := range newRoom.Clients {
			client.Write(fmt.Sprintf("%s entered the room", au.Green(c.Name)))
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
	c.Write(fmt.Sprintf("what is your %s?", au.Magenta("name")))
	data, err := c.reader.ReadString('\n')
	if err != nil {
		return err
	}
	c.Name = strings.TrimSpace(string(data))
	c.Write(fmt.Sprintf("%s, %s\r\n", au.Framed("Welcome"), au.Bold(c.Name)))
	c.messages <- NewMessage(ClientStartedMessage, c, "login completed", []string{})
	return nil
}

// closeConnection closes the client connection connection
func (c *Client) closeConnection(reason string) {
	c.conn.Close()
	c.messages <- NewMessage(ClientStoppedMessage, c, reason, []string{})
}

// handleConnection handles a new network client connection
func (c *Client) handleConnection() {
	err := c.handleLogin()
	if err != nil {
		c.closeConnection("login failed")
		return
	}

	c.Write(c.room.Look())

	for {
		data, err := c.reader.ReadString('\n')
		if err != nil {
			c.closeConnection("read failed")
			return
		}

		message := strings.TrimSpace(string(data))
		fields := strings.Fields(message)
		if len(fields) < 1 {
			continue
		}
		command := fields[0]
		args := []string{}
		target := ""

		if len(fields) > 1 {
			args = fields[1:]
			target = args[0]
		}

		switch command {
		case "exit":
			c.closeConnection("client exited")
			return
		case "look":
			c.messages <- NewMessage(ClientLookMessage, c, target, args)
		case "enter":
			c.messages <- NewMessage(ClientEnterMessage, c, target, args)
		case "say":
			msg := strings.Join(args, " ")
			c.messages <- NewMessage(ClientChatMessage, c, msg, args)
		default:
			msg := fmt.Sprintf("unknown command:%s, args:%v", command, args)
			c.messages <- NewMessage(ErrorMessage, c, msg, args)
		}
	}
}
