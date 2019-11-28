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
	defer func() { c.log("handleInput stopped") }()
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
		c.log("pub start input: %s", input)
		c.Input <- NewClientInputMessage(c, input)
		c.log("pub done input: %s", input)
	}
}
