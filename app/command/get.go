package command

import (
	"errors"
)

func (cmd Command) Get(key interface{}) error {
	parsedKey, ok := key.(string)
	if !ok {
		return errors.New("key is not a string")
	}

	if value, ok := cmd.client.Cache.Get(parsedKey); ok {
		return cmd.client.Send(value)
	}

	return cmd.client.Send(nil)
}
