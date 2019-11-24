package main

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"

	au "github.com/logrusorgru/aurora"
)

// Client is a client
type Client struct {
	ID       string
	IsClosed bool
	Upstream chan<- interface{}
	conn     net.Conn
	reader   *bufio.Reader
	Player   *Player
}

// NewClient creates a client
func NewClient(conn net.Conn, upstream chan<- interface{}) *Client {
	reader := bufio.NewReader(conn)
	client := &Client{GetID(), false, upstream, conn, reader, nil}

	go func() {
		s := "NewClient | getPlayer"
		defer Trace(s, "ended")
		Trace(s, "started")
		err := client.handleLogin()
		if err != nil {
			Trace(s, "login failed")
			client.closeConnection("login failed")
		} else {
			Trace(s, "login succeeded")
			client.Upstream <- NewClientJoinMessage(client)
			go client.handleInput()
		}
	}()

	return client
}

// Write puts a message in the writeq channel:
func (c *Client) Write(message string, m ...interface{}) {
	msg := fmt.Sprintf(message, m...)
	c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	c.conn.Write([]byte(msg + "\r\n"))
}

// Read blocks to read a line of input from the socket with a timeout of 1s
func (c *Client) Read() string {
	readTimeout := time.Second * 5
	c.conn.SetReadDeadline(time.Now().Add(readTimeout))
	data, err := c.reader.ReadString('\n')
	if err != nil {
		if e, ok := err.(net.Error); ok && e.Timeout() {
			return ""
		}
		msg := fmt.Sprintf("read failed: %s", err)
		c.closeConnection(msg)
		return ""
	}
	return data
}

// closeConnection closes the client connection connection
func (c *Client) closeConnection(reason string) {
	s := fmt.Sprintf("c | %s | closeConnection", c.ID)
	defer Trace(s, "ended")
	Trace(s, "started")

	if c.IsClosed {
		Trace(s, "already closed")
		return
	}
	c.IsClosed = true
	c.conn.Close()
	msg := NewClientLeaveMessage(c, reason)
	c.Upstream <- msg
}

// handleLogin handles the client login
func (c *Client) handleLogin() error {
	s := fmt.Sprintf("c | %s | handleLogin", c.ID)
	defer Trace(s, "ended")
	Trace(s, "started")

	// ask for name:
	msg := fmt.Sprintf("what is your %s? ", au.Magenta("name"))
	c.conn.Write([]byte(msg))

	// read name:
	data, err := c.reader.ReadString('\n')
	if err != nil {
		Trace(s, "failed to read name. error: %s", err)
		return err
	}
	if len(data) < 1 {
		Trace(s, "failed to provide name")
		return fmt.Errorf("failed to provide name")
	}

	// get player:
	name := strings.TrimSpace(string(data))
	player := NewPlayer(c, name)
	c.Player = player
	c.Write("Welcome, %s\r\n", au.Bold(c.Player.Name))
	return nil
}

// handleInput handles input from the network connection
func (c *Client) handleInput() {
	s := fmt.Sprintf("c | %s | handleInput", c.ID)
	for {
		if c.IsClosed {
			Trace(s, "closing")
			return
		}

		data := c.Read()
		input := strings.TrimSpace(string(data))
		if len(input) < 1 {
			continue
		}

		inputMsg := NewInputMessage(c, input)
		c.Upstream <- inputMsg
		Trace(s, "sent upstream InputMessage id:%s, input:%s", inputMsg.Meta.ID, inputMsg.Input)

		if input == "exit" {
			Trace(s, "exiting")
			c.closeConnection("client exited")
			return
		}
	}
}
