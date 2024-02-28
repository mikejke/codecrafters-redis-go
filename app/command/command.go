package command

import "github.com/codecrafters-io/redis-starter-go/app/client"

type Command struct {
	client *client.Client
}

func NewCommand(client *client.Client) *Command {
	return &Command{
		client: client,
	}
}
