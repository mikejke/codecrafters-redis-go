package command

func (cmd Command) Echo(content []interface{}) {
	cmd.client.Send(content...)
}
