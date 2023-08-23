package horus

import (
	"errors"
	"fmt"
)

var (
	ErrExist    = errors.New("already exist")
	ErrNotExist = errors.New("not exist")
)

type InvalidError struct {
	What string
	Why  string
}

func (e *InvalidError) Error() string {
	return fmt.Sprintf(`invalid value of "%s": %s`, e.What, e.Why)
}
