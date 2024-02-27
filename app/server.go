package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(conn)
	}
}

const (
	PING = "PING"
	ECHO = "ECHO"
	SET  = "SET"
	GET  = "GET"
)

func handleConnection(conn net.Conn) {
	client, err := NewClient(conn)
	if err != nil {
		fmt.Println("Error: ", err.Error())
		return
	}
	defer client.Close()

	for {
		result, err := client.Read()
		if err == io.EOF {
			fmt.Println("EOF")
			return
		}

		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}

		fmt.Println("RECIEVED", result)

		switch c := result.Content().(type) {
		case []interface{}:
			switch strings.ToUpper(fmt.Sprintf("%v", c[0])) {
			case PING:
				client.Send([]interface{}{"PONG"})
			case ECHO:
				client.Send(c[1:])
			case SET:
				key, ok := c[1].(string)
				if !ok {
					fmt.Println("key is not a string")
					return
				}

				value := c[2]
				client.Store(key, value)
				client.Send([]interface{}{"OK"})
			case GET:
				key, ok := c[1].(string)
				if !ok {
					fmt.Println("key is not a string")
					return
				}
				storedValue := client.Get(key)
				client.Send(storedValue)
			}
		default:
			fmt.Println("unknown command")
			return
		}
	}
}
