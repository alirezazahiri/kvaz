package command

import "kvaz/internal/storage"

type ExpireCommand struct{ Key string }

func NewExpireCommand(k string) ExpireCommand {
	return ExpireCommand{Key: k}
}
func (c ExpireCommand) cmd()                            {}
func (c ExpireCommand) Execute(_ storage.StorageEngine) {}
