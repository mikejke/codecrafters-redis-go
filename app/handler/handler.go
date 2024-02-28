package handler

import (
	"fmt"
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

func HandleCommand(client *client.Client, content []interface{}) {
	cmd := command.NewCommand(client)

	switch strings.ToUpper(fmt.Sprintf("%v", content[0])) {
	case ping:
		cmd.Ping()
	case echo:
		cmd.Echo(content[1:])
	case set:
		cmd.Set(content[1:])
	case get:
		cmd.Get(content[1])
	case config:
		cmd.Config(content[1:])
	}
}
