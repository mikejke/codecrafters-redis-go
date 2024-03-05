package command

func (cmd Command) Ping() error {
	return cmd.client.Send("PONG")
}
