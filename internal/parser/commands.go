package parser

import "kvaz/internal/command"

func Ping() (command.PingCommand, error) {
	return command.NewPingCommand(), nil
}

func Get(k string) (command.GetCommand, error) {
	return command.NewGetCommand(k), nil
}

func Set(k, v string) (command.SetCommand, error) {
	return command.NewSetCommand(k, v), nil
}

func Del(k string) (command.DelCommand, error) {
	return command.NewDelCommand(k), nil
}

func Exists(k string) (command.ExistsCommand, error) {
	return command.NewExistsCommand(k), nil
}

func Expire(k string, secs int64) (command.ExpireCommand, error) {
	return command.NewExpireCommand(k, secs), nil
}

func TTL(k string) (command.TTLCommand, error) {
	return command.NewTTLCommand(k), nil
}
