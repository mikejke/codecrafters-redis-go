package command

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/codecrafters-io/redis-starter-go/app/cache"
)

const PX = "PX"

func (cmd Command) Set(args []interface{}) error {
	key, ok := args[0].(string)
	if !ok {
		return errors.New("key is not a string")
	}

	value := args[1]
	item := &cache.Item{
		Key:   key,
		Value: value,
	}

	// TODO: refactor this block of code
	if len(args[0:]) == 4 && strings.ToUpper(fmt.Sprintf("%v", args[2])) == PX {
		expiryTime, ok := args[3].(string)
		if !ok {
			return errors.New("expiry time is not integer")
		}

		parsedTime, err := strconv.Atoi(expiryTime)
		if err != nil {
			return errors.New("failed to parse TTL")
		}

		item.TTL = time.Duration(parsedTime)
	}

	cmd.client.Cache.Set(item)

	return cmd.client.Send("OK")
}
