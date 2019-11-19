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
	// login and exit if fails:
	err := c.handleLogin()
	if err != nil {
		c.closeConnection("login failed")
		return
	}

	// display room once after succesful login:
	c.Write(c.room.Look(c.Prompt()))

	for {
		// wait for a message:
		data, err := c.reader.ReadString('\n')
		if err != nil {
			c.closeConnection("read failed")
			return
		}

		// split the message up into fields:
		message := strings.TrimSpace(string(data))
		fields := strings.Fields(message)

		// nothing to do if empty:
		if len(fields) < 1 {
			continue
		}

		// split the fields into 3 parts:
		// * command (first keyword)
		// * args    (list of all keywords after command)
		// * target  (second keyword or empty string if not exists)
		//   - target is syntactic sugar for commands that require/use 2nd parm
		// examples:
		// "open door" = (command:open, args:[door], target:door)
		// "look"      = (command:look, args:[],     target:"")
		command := fields[0]
		args := []string{}
		target := ""

		// only set args and target if there are more keywords:
		if len(fields) > 1 {
			args = fields[1:]
			target = args[0]
		}

		// handle the command and exit if appropriate:
		isExiting := c.handleCommand(command, target, args)
		if isExiting {
			c.closeConnection("client exited")
			return
		}
	}
}

// handleCommand will handle any command entered by a client
func (c *Client) handleCommand(command, target string, args []string) bool {
	isExiting := false
	switch command {
	case "exit":
		isExiting = true
	case "look":
		c.messages <- NewMessage(ClientLookMessage, c, target, args)
	case "enter":
		c.messages <- NewMessage(ClientEnterMessage, c, target, args)
	case "say":
		msg := strings.Join(args, " ")
		c.messages <- NewMessage(ClientChatMessage, c, msg, args)
	case "debug":
		c.InCombat = !c.InCombat
	default:
		msg := fmt.Sprintf("unknown command:%s, args:%v", command, args)
		c.messages <- NewMessage(ErrorMessage, c, msg, args)
	}
	return isExiting
}
