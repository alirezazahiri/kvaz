package command

import (
	"fmt"
	"kvaz/internal/storage"
)

type DelCommand struct{ Key string }

func NewDelCommand(k string) DelCommand {
	return DelCommand{Key: k}
}
func (c DelCommand) cmd() {}
func (c DelCommand) Execute(_ storage.StorageEngine) {
	fmt.Println("DelCommand")
}
