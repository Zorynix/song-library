package errors

import "errors"

var (
	ErrNotFound        = errors.New("resource not found")
	ErrInternal        = errors.New("internal server error")
	ErrBadRequest      = errors.New("bad request")
	ErrInvalidInput    = errors.New("invalid input data")
	ErrOperationFailed = errors.New("operation failed")
)
