package client

import (
	"fmt"
	"io"
	"net"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
)

var (
	_ io.Closer = &Client{}
)

type Client struct {
	conn   net.Conn
	reader *Reader
	writer *Writer
	Cache  *cache.Cache
}

func NewClient(conn net.Conn) (*Client, error) {
	return &Client{
		conn:   conn,
		reader: NewReader(conn),
		writer: NewWriter(conn),
		Cache:  cache.NewCache(),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Read() (*Result, error) {
	return c.reader.Read()
}

func (c *Client) Send(values ...interface{}) error {
	if err := c.writer.WriteArray(values); err != nil {
		return fmt.Errorf("failed to execute operation: %v", values[0])
	}

	return nil
}
