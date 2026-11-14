package command

import (
	"fmt"
	"kvaz/internal/storage"
)

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
func (c SetCommand) Execute(storage storage.StorageEngine) {
	ok := storage.Set(c.Key, c.Value)
	if ok {
		fmt.Println("OK")
		return
	}
	fmt.Println("ERROR")
}
