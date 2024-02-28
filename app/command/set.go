package command

import (
	"fmt"
	"strconv"
	"strings"
)

const PX = "PX"

func (cmd Command) Set(args []interface{}) {
	key, ok := args[0].(string)
	if !ok {
		fmt.Println("key is not a string")
		return
	}

	value := args[1]
	cmd.client.Store(key, value)

	// TODO: refactor this block of code
	if len(args[0:]) == 4 && strings.ToUpper(fmt.Sprintf("%v", args[2])) == PX {
		expiryTime, ok := args[3].(string)
		if !ok {
			fmt.Println("expiry time is not integer")
			return
		}
		fmt.Println(expiryTime)

		parsedTime, err := strconv.Atoi(expiryTime)
		if err != nil {
			fmt.Println("failed to parse expiry time")
			return
		}

		cmd.client.SetExpirationTime(key, parsedTime)
	}

	cmd.client.Send("OK")
}
