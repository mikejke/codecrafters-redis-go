package main

import (
	"fmt"
	"io"
	"net"
	"time"
)

var (
	_ io.Closer = &Client{}
)

type Client struct {
	conn    net.Conn
	reader  *Reader
	writer  *Writer
	storage map[string]interface{}
}

func NewClient(conn net.Conn) (*Client, error) {
	return &Client{
		conn:    conn,
		reader:  NewReader(conn),
		writer:  NewWriter(conn),
		storage: make(map[string]interface{}),
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

func (c *Client) Store(key string, value interface{}) {
	c.storage[key] = value
}

func (c *Client) Get(key string) interface{} {
	if value, ok := c.storage[key]; ok {
		return value
	}

	return nil
}

func (c *Client) SetExpirationTime(key string, expirationTime int) {
	fmt.Println("planned to delete", key)
	go func() {
		time.Sleep(time.Millisecond * time.Duration(expirationTime))
		delete(c.storage, key)
		fmt.Println("key", key, "was deleted")
	}()
}
