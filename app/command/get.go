package command

import "fmt"

func (cmd *Command) Get(key interface{}) {
	parsedKey, ok := key.(string)
	if !ok {
		fmt.Println("key is not a string")
		return
	}
	cmd.client.Send(cmd.client.Get(parsedKey))
}
