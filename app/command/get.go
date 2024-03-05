package command

import (
	"errors"
)

func (cmd Command) Get(key interface{}) error {
	parsedKey, ok := key.(string)
	if !ok {
		return errors.New("key is not a string")
	}
	return cmd.client.Send(cmd.client.Cache.Get(parsedKey))
}
