package command

type ExpireCommand struct{ Key string }

func NewExpireCommand(k string) ExpireCommand {
	return ExpireCommand{Key: k}
}
func (c ExpireCommand) cmd() {}
