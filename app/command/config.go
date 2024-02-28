package command

import (
	"fmt"
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
				return fmt.Errorf("key is not a string")
			}

			cmd.client.Send(config.Config.Get(key))
		default:
			return fmt.Errorf("unknown subcommand")
		}
	default:
		return fmt.Errorf("unknown subcommand")
	}
	return nil
}
