package command

import "fmt"

type PingCommand struct{}

func NewPingCommand() PingCommand {
	return PingCommand{}
}
func (c PingCommand) cmd() {}
func (c PingCommand) Execute() {
	fmt.Println("PingCommand")
}
