// Package processor: this one is responsible for processing commands
package processor

import (
	"fmt"
	"kvaz/internal/command"
	"kvaz/internal/dto"
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

func (p *Processor) Process(cmd command.Command) dto.Response {
	fmt.Printf("executing task: %T\n", cmd)
	switch cmd := cmd.(type) {
	case command.PingCommand:
		return dto.EmptySuccessResponse()
	case command.SetCommand:
		ok := p.storage.Set(cmd.Key, cmd.Value)
		if ok {
			return dto.EmptySuccessResponse()
		}
		return dto.ErrorResponse(fmt.Errorf("failed to set value for key '%s'", cmd.Key))
	case command.GetCommand:
		value, err := p.storage.Get(cmd.Key)
		if err != nil {
			return dto.ErrorResponse(err)
		}
		return dto.SuccessResponse(value)
	case command.DelCommand:
		ok := p.storage.Delete(cmd.Key)
		if ok {
			return dto.EmptySuccessResponse()
		}
		return dto.ErrorResponse(fmt.Errorf("failed to delete key '%s'", cmd.Key))
	case command.ExistsCommand:
		ok := p.storage.Exists(cmd.Key)
		if ok {
			return dto.EmptySuccessResponse()
		}
		return dto.NotFoundResponse()
	// case command.ExpireCommand:
	// case command.TTLCommand:
	default:
		todo()
		return dto.Response{}
	}
}

func todo() {
	fmt.Println("not implemented yet!")
}
