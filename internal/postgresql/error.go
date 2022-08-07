package postgresql

import (
	"errors"
)

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("postgresql: already exists")
	ErrNotExist = errors.New("postgresql: does not exist")

	// storage errors
	ErrRetry    = errors.New("postgresql: retry storage operation")
	ErrConflict = errors.New("postgresql: conflict in storage operation")
)
