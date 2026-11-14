package command

type DelCommand struct{ Key string }

func NewDelCommand(k string) DelCommand {
	return DelCommand{Key: k}
}
func (c DelCommand) cmd() {}
