package storage

type RecordState uint8

const (
	RecordStateMissing    = 1
	RecordStatePersistent = 2
	RecordStateExpiring   = 3
)
