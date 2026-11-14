// Package command: this one is responsible for command objects and their payloads
package command

type Command interface {
	cmd()
}
