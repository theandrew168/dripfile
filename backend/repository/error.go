package repository

import "errors"

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("repository: already exists")
	ErrNotExist = errors.New("repository: does not exist")

	// storage errors
	ErrRetry    = errors.New("repository: retry storage operation")
	ErrConflict = errors.New("repository: conflict in storage operation")
)
