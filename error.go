package goregexpopulate

import "errors"

var (
	ErrNonPointer   error = errors.New("parameter must be a pointer")
	ErrNilPointer   error = errors.New("parameter must not be nil")
	ErrEmptyPattern error = errors.New("pattern is empty")
)
