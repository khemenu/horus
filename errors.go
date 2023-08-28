package horus

import (
	"errors"
)

var (
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrNotExist           = errors.New("not exist")
	ErrExist              = errors.New("already exist")
	ErrFailedPrecondition = errors.New("failed precondition")
)
