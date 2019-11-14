package network

import (
	"bufio"
	"net"
)

// Client is a network client
type Client struct {
	Addr   string
	Conn   net.Conn
	Reader *bufio.Reader
}

// NewClient creates a network client
func NewClient(conn net.Conn) *Client {
	addr := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	return &Client{Addr: addr, Conn: conn, Reader: reader}
}

// Write will write a message to a network client
func (c *Client) Write(message string) {
	c.Conn.Write([]byte(message))
}
