package command

func (cmd Command) Echo(content []interface{}) error {
	return cmd.client.Send(content...)
}
