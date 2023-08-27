package horus

import (
	"errors"
)

var (
	ErrExist    = errors.New("already exist")
	ErrNotExist = errors.New("not exist")
)
