package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

// Client is a client
type Client struct {
	ID       string
	Input    chan<- *ClientInputMessage
	IsClosed bool
	conn     net.Conn
	reader   *bufio.Reader
}

// NewClient creates a client
func NewClient(connection net.Conn, input chan<- *ClientInputMessage) *Client {
	reader := bufio.NewReader(connection)
	client := &Client{NewID(), input, false, connection, reader}
	go client.handleInput()
	return client
}

// Write writes a message to the client
func (c *Client) Write(message string, args ...interface{}) {
	msg := fmt.Sprintf(message, args...)
	c.conn.Write([]byte(fmt.Sprintf("%s\r\n", msg)))
}

// Close closes the client connection
func (c *Client) Close() {
	c.IsClosed = true
	c.conn.Close()
	c.log("closed connection")
}

// log is for logging a message
func (c *Client) log(message string, args ...interface{}) {
	msg := fmt.Sprintf("client:%s | %s\n", c.ID, message)
	log.Printf(msg, args...)
}

// handleInput handles client input:
func (c *Client) handleInput() {
	for {
		if c.IsClosed {
			c.log("handleInput stopping per closed connection")
			return
		}
		input, err := c.reader.ReadString('\n')
		if err != nil {
			c.log("handleInput stopping per reader error: %+v", err)
			c.Close()
			return
		}
		c.Input <- NewClientInputMessage(c, input)
	}
}
