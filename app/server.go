package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		defer conn.Close()
		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Error while reading: ", err.Error())
			continue
		}

		fmt.Printf("Received message: %s", string(buf[:n]))
		if _, err := conn.Write([]byte("+PONG\r\n")); err != nil {
			fmt.Println("Error while responding: ", err.Error())
			return
		}
	}

}
