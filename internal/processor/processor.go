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
		return dto.SimpleString("PONG")
	case command.SetCommand:
		ok := p.storage.Set(cmd.Key, cmd.Value)
		if ok {
			return dto.SimpleString("OK")
		}
		return dto.Error(fmt.Sprintf("failed to set value for key '%s'", cmd.Key))
	case command.GetCommand:
		value, err := p.storage.Get(cmd.Key)
		if err != nil {
			return dto.Null()
		}
		return dto.BulkString(value)
	case command.DelCommand:
		ok := p.storage.Delete(cmd.Key)
		if ok {
			return dto.Integer(1)
		}
		return dto.Integer(0)
	case command.ExistsCommand:
		ok := p.storage.Exists(cmd.Key)
		if ok {
			return dto.Integer(1)
		}
		return dto.Integer(0)
	case command.TTLCommand:
		secs, state := p.storage.TTL(cmd.Key)
		switch state {
		case storage.RecordStatePersistent:
			return dto.Integer(-1)
		case storage.RecordStateMissing:
			return dto.Integer(-2)
		case storage.RecordStateExpiring:
			return dto.Integer(secs)
		default:
			// This obviously doesn't happen, but in case we add some new state I didn't treat the remaining case as default
			return dto.Null()
		}
	case command.ExpireCommand:
		ok := p.storage.Expire(cmd.Key, cmd.Seconds)
		if ok {
			return dto.Integer(1)
		}
		return dto.Integer(0)
	default:
		return dto.Error("unknown command")
	}
}
