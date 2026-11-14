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

func (p *Processor) Process(cmd command.Command) string {
	fmt.Printf("executing task: %T\n", cmd)
	switch cmd := cmd.(type) {
	case command.PingCommand:
		return "PONG"
	case command.SetCommand:
		ok := p.storage.Set(cmd.Key, cmd.Value)
		if ok {
			return "OK"
		}
		return "ERROR"
	case command.GetCommand:
		value, err := p.storage.Get(cmd.Key)
		if err != nil {
			fmt.Println(err.Error())
			return "ERROR"
		}
		return value
	// case command.DelCommand:
	// case command.ExpireCommand:
	// case command.ExistsCommand:
	// case command.TTLCommand:
	default:
		todo()
		return ""
	}
}

func todo() {
	fmt.Println("not implemented yet!")
}

// TODO: introduce some type for each command's output
