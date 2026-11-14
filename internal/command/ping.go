package command

type PingCommand struct{}

func NewPingCommand() PingCommand {
	return PingCommand{}
}
func (c PingCommand) cmd() {}
