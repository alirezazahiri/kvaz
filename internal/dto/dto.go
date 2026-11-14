// Package dto: this one is responsible for the construction of data-transfer-object
package dto

type Response struct {
	Status Status
	Value  any
	Error  error
}

func EmptySuccessResponse() Response {
	return Response{
		Status: StatusOk,
		Value:  nil,
		Error:  nil,
	}
}

func SuccessResponse(v any) Response {
	return Response{
		Status: StatusOk,
		Value:  v,
		Error:  nil,
	}
}

func NotFoundResponse() Response {
	return Response{
		Status: StatusNotFound,
		Value:  nil,
		Error:  nil,
	}
}

func ErrorResponse(err error) Response {
	return Response{
		Status: StatusError,
		Value:  nil,
		Error:  err,
	}
}
