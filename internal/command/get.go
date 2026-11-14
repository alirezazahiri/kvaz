package command

import (
	"fmt"
	"kvaz/internal/storage"
)

type GetCommand struct{ Key string }

func NewGetCommand(k string) GetCommand {
	return GetCommand{Key: k}
}
func (c GetCommand) cmd() {}
func (c GetCommand) Execute(s storage.StorageEngine) {
	value, err := s.Get(c.Key)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(value)
}
