package command

type TTLCommand struct {
	Key string
}

func NewTTLCommand(k string) TTLCommand {
	return TTLCommand{
		Key: k,
	}
}
func (c TTLCommand) cmd() {}
