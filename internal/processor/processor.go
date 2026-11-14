// Package processor: this one is responsible for processing commands
package processor

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/storage"
)

type Processor struct {
	storage storage.StorageEngine
}

func New(s storage.StorageEngine) *Processor {
	return &Processor{
		storage: s,
	}
}

func (p *Processor) Process(cmd command.Command) {
	fmt.Printf("executing task: %#v\n", cmd)
	cmd.Execute(p.storage)
	// TODO: execute commands and treat them like objects only
}
