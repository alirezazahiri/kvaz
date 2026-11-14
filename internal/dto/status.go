package dto

type Status string

const (
	StatusOk       Status = "OK"
	StatusNotFound Status = "NOT_FOUND"
	StatusError    Status = "ERROR"
)
