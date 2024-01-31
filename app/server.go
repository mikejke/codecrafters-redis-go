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
			fmt.Printf("Client[%s] ended connection\n", client.conn.RemoteAddr().String())
			return
		}

		if err != nil {
			fmt.Println("Error: ", err.Error())
			return
		}

		content := result.Content()

		fmt.Printf("Received message from client [%s]: %v\n", client.conn.RemoteAddr().String(), content)

		switch c := content.(type) {
		case string:
			if strings.ToUpper(c) == PING {
				fmt.Println("Recieved ping (1) command, responeded with '+PONG\r\n'")
				client.conn.Write([]byte("+PONG\r\n"))
			}
		case []interface{}:
			switch strings.ToUpper(fmt.Sprintf("%v", c[0])) {
			case PING:
				fmt.Println("Recieved ping (2) command, responeded with '+PONG\r\n'")
				client.conn.Write([]byte("+PONG\r\n"))
			case ECHO:
				var message string
				for i := 1; i < len(c); i++ {
					message += fmt.Sprintf("%v ", c[i])
				}
				fmt.Println("Recieved echo command, responeded with '", strings.TrimRight(message, " "), "'")
				client.conn.Write([]byte(strings.TrimRight(message, " ")))
			}
		default:
			fmt.Println("unknown command")
			return
		}
	}
}
