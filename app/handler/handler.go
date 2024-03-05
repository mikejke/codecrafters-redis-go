package handler

import (
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/client"
	"github.com/codecrafters-io/redis-starter-go/app/command"
)

const (
	ping   = "PING"
	echo   = "ECHO"
	set    = "SET"
	get    = "GET"
	config = "CONFIG"
)

func HandleCommand(client *client.Client, content []interface{}) error {
	cmd := command.NewCommand(client)

	if command, ok := content[0].(string); ok {
		switch strings.ToUpper(command) {
		case ping:
			return cmd.Ping()
		case echo:
			return cmd.Echo(content[1:])
		case set:
			return cmd.Set(content[1:])
		case get:
			return cmd.Get(content[1])
		case config:
			return cmd.Config(content[1:])
		}
	}

	return nil
}
