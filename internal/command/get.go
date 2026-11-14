package command

type GetCommand struct{ Key string }

func NewGetCommand(k string) GetCommand {
	return GetCommand{Key: k}
}
func (c GetCommand) cmd() {}
