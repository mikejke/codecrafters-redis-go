package main

import (
	"io"
	"net"
)

var (
	_ io.Closer = &Client{}
)

type Client struct {
	conn   net.Conn
	reader *Reader
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Read() (*Result, error) {
	return c.reader.Read()
}

func NewClient(conn net.Conn) (*Client, error) {
	return &Client{
		conn:   conn,
		reader: NewReader(conn),
	}, nil
}
