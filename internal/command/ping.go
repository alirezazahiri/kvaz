package command

import (
	"fmt"
	"kvaz/internal/storage"
)

type PingCommand struct{}

func NewPingCommand() PingCommand {
	return PingCommand{}
}
func (c PingCommand) cmd() {}
func (c PingCommand) Execute(_ storage.StorageEngine) {
	fmt.Println("PingCommand")
}
