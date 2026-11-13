package command

import "fmt"

type SetCommand struct {
	Key   string
	Value string
}

func NewSetCommand(k string, v string) SetCommand {
	return SetCommand{
		Key:   k,
		Value: v,
	}
}
func (c SetCommand) cmd() {}
func (c SetCommand) Execute() {
	fmt.Println("SetCommand")
}
