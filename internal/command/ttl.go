package command

type TTLCommand struct{}

func NewTTLCommand(k string) TTLCommand {
	return TTLCommand{}
}
func (c TTLCommand) cmd()     {}
func (c TTLCommand) Execute() {}
