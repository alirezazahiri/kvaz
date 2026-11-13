package command

import "fmt"

type DelCommand struct{ Key string }

func NewDelCommand(k string) DelCommand {
	return DelCommand{Key: k}
}
func (c DelCommand) cmd() {}
func (c DelCommand) Execute() {
	fmt.Println("DelCommand")
}
