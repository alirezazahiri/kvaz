package command

type ExistsCommand struct{ Key string }

func NewExistsCommand(k string) ExistsCommand {
	return ExistsCommand{Key: k}
}
func (c ExistsCommand) cmd()     {}
func (c ExistsCommand) Execute() {}
