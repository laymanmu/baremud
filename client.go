package main

import (
	"bufio"
	"fmt"
	"net"
)

// Client is a client
type Client struct {
	ID       string
	Input    chan<- *ClientInputMessage
	IsClosed bool
	conn     net.Conn
	reader   *bufio.Reader
	log      Logger
}

// NewClient creates a client
func NewClient(conn net.Conn, input chan<- *ClientInputMessage) *Client {
	id := NewID("client")
	log := NewLogger(id)
	reader := bufio.NewReader(conn)
	client := &Client{id, input, false, conn, reader, log}
	go client.handleInput()
	return client
}

// Write writes a message to the client
func (c *Client) Write(message string, args ...interface{}) {
	defer (Track("Write", c.log))()
	msg := fmt.Sprintf(message, args...)
	c.conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
}

// Close closes the client connection
func (c *Client) Close() {
	defer (Track("Close", c.log))()
	c.IsClosed = true
	c.conn.Close()
}

// handleInput handles client input:
func (c *Client) handleInput() {
	defer (Track("handleInput", c.log))()
	for {
		if c.IsClosed {
			return
		}
		input, err := c.reader.ReadString('\n')
		if err != nil {
			if c.IsClosed {
				return
			}
			c.log("closing from error: %s", err.Error())
			c.Close()
			return
		}
		c.log("sending input message: %s", input)
		c.Input <- NewClientInputMessage(c, input)
	}
}
