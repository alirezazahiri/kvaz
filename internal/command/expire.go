package command

type ExpireCommand struct {
	Key     string
	Seconds int64
}

func NewExpireCommand(k string, secs int64) ExpireCommand {
	return ExpireCommand{Key: k, Seconds: secs}
}
func (c ExpireCommand) cmd() {}
