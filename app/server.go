package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
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
	PX   = "PX"
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

		switch c := result.Content().(type) {
		case []interface{}:
			switch strings.ToUpper(fmt.Sprintf("%v", c[0])) {
			case PING:
				client.Send("PONG")
			case ECHO:
				client.Send(c[1:]...)
			case SET:
				key, ok := c[1].(string)
				if !ok {
					fmt.Println("key is not a string")
					return
				}

				value := c[2]
				client.Store(key, value)

				if len(c[1:]) == 4 && strings.ToUpper(fmt.Sprintf("%v", c[3])) == PX {
					expiryTime, ok := c[4].(string)
					if !ok {
						fmt.Println("expiry time is not integer")
						return
					}
					fmt.Println(expiryTime)

					parsedTime, err := strconv.Atoi(expiryTime)
					if err != nil {
						fmt.Println("failed to parse expiry time")
						return
					}

					client.SetExpirationTime(key, parsedTime)
				}

				client.Send("OK")
			case GET:
				key, ok := c[1].(string)
				if !ok {
					fmt.Println("key is not a string")
					return
				}
				client.Send(client.Get(key))
			}
		default:
			fmt.Println("unknown command")
			return
		}
	}
}
