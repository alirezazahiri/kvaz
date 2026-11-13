package command

import "fmt"

type GetCommand struct{ Key string }

func NewGetCommand(k string) GetCommand {
	return GetCommand{Key: k}
}
func (c GetCommand) cmd() {}
func (c GetCommand) Execute() {
	fmt.Println("GetCommand")
}
