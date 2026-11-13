package parser

import "kvaz/internal/command"

func (p *Parser) ping() (command.PingCommand, error) {
	return command.NewPingCommand(), nil
}

func (p *Parser) get() (command.GetCommand, error) {
	return command.NewGetCommand(""), nil
}

func (p *Parser) set() (command.SetCommand, error) {
	return command.NewSetCommand("", ""), nil
}

func (p *Parser) del() (command.DelCommand, error) {
	return command.NewDelCommand(""), nil
}

func (p *Parser) exists() (command.ExistsCommand, error) {
	return command.NewExistsCommand(""), nil
}

func (p *Parser) expire() (command.ExpireCommand, error) {
	return command.NewExpireCommand(""), nil
}

func (p *Parser) ttl() (command.TTLCommand, error) {
	return command.NewTTLCommand(""), nil
}
