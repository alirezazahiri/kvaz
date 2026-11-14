// Package command: this one is responsible for command objects and their payloads
package command

import "kvaz/internal/storage"

type Command interface {
	cmd()
	Execute(s storage.StorageEngine)
}
