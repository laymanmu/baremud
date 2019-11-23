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
	Health   int
	InCombat bool
	room     *Room
	conn     net.Conn
	reader   *bufio.Reader
	messages chan<- *Message
}

// NewClient creates a client
func NewClient(conn net.Conn, messages chan<- *Message, room *Room) *Client {
	name := ""
	health := 100
	inCombat := false
	reader := bufio.NewReader(conn)
	client := &Client{GetID(), name, health, inCombat, room, conn, reader, messages}
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
		c.Write(c.room.Look(c.Prompt()))
	}
}

// Prompt sends the prompt to the client
func (c *Client) Prompt() string {
	left := au.BrightBlack("-=[")
	right := au.BrightBlack("]=-")
	div := au.BrightBlack("|")
	hp := au.Red(fmt.Sprintf("%v", c.Health))
	alert := au.Green("ok")
	if c.InCombat {
		alert = au.Red("combat")
	}
	if c.Health < 100 {
		return fmt.Sprintf("%s%s%s%s%s", left, hp, div, alert, right)
	}
	return fmt.Sprintf("%s%s%s", left, alert, right)
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
	return nil
}

// closeConnection closes the client connection connection
// messages: ClientStoppedMessage
func (c *Client) closeConnection(reason string) {
	c.conn.Close()
	c.messages <- NewMessage(ClientStoppedMessage, c, reason, []string{})
}

// handleConnection handles a new network client connection
// messages: ClientStartedMessage, ClientInputMessage
func (c *Client) handleConnection() {
	// login and exit if fails:
	err := c.handleLogin()
	if err != nil {
		c.closeConnection("login failed")
		return
	}
	// publish a start message:
	c.messages <- NewMessage(ClientStartedMessage, c, "login completed", []string{})

	// read from the client and publish input messages:
	for {
		data, err := c.reader.ReadString('\n')
		if err != nil {
			c.closeConnection("read failed")
			return
		}
		input := strings.TrimSpace(string(data))
		args  := strings.Fields(input)
		if len(args) < 1 {
			continue
		}
		c.messages <- NewMessage(InputMessage, c, input, args)
	}
}
