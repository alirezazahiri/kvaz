package parser

import "kvaz/internal/command"

func (p *Parser) ping() (command.PingCommand, error) {
	return command.NewPingCommand(), nil
}

func (p *Parser) get(k string) (command.GetCommand, error) {
	return command.NewGetCommand(k), nil
}

func (p *Parser) set(k, v string) (command.SetCommand, error) {
	return command.NewSetCommand(k, v), nil
}

func (p *Parser) del(k string) (command.DelCommand, error) {
	return command.NewDelCommand(k), nil
}

func (p *Parser) exists(k string) (command.ExistsCommand, error) {
	return command.NewExistsCommand(k), nil
}

func (p *Parser) expire(k string) (command.ExpireCommand, error) {
	return command.NewExpireCommand(k), nil
}

func (p *Parser) ttl(k string) (command.TTLCommand, error) {
	return command.NewTTLCommand(k), nil
}
