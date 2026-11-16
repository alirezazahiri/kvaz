package resp2parser

type CommandKind uint8

const (
	CommandKindGet CommandKind = iota + 1
	CommandKindSet
	CommandKindPing
	CommandKindDel
	CommandKindExpire
	CommandKindExists
	CommandKindTTL
)
