package command

func (cmd *Command) Ping() {
	cmd.client.Send("PONG")
}
