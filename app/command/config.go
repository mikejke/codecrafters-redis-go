package command

import (
	"errors"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/app/config"
)

var (
	GET = "GET"
)

func (cmd Command) Config(args []interface{}) error {
	switch sub := args[0].(type) {
	case string:
		switch strings.ToUpper(sub) {
		case GET:
			key, ok := args[1].(string)
			if !ok {
				return errors.New("key is not a string")
			}

			cmd.client.Send(key, config.Config.Get(key))
		default:
			return errors.New("unknown subcommand")
		}
	default:
		return errors.New("unknown subcommand")
	}
	return nil
}
