package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/config"
	"github.com/codecrafters-io/redis-starter-go/app/handler"
)

func main() {
	fmt.Println("Logs from your program will appear here!")

	flag.StringVar(&config.Config.Dir, "dir", "", "The directory where RDB files are stored")
	flag.StringVar(&config.Config.RDBFilename, "dbfilename", "", "The name of the RDB file")
	flag.Parse()

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

func handleConnection(conn net.Conn) {
	client, err := client.NewClient(conn)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer client.Close()

	for {
		result, err := client.Read()
		if err == io.EOF {
			fmt.Println("Error: EOF")
			return
		}

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		if content, ok := result.Content().([]interface{}); ok {
			err := handler.HandleCommand(client, content)
			if err != nil {
				fmt.Println("unexpected content")
				return
			}
		} else {
			fmt.Println("unexpected content", content)
			return
		}
	}
}
