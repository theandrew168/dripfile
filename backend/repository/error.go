package repository

import "errors"

var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("repository: already exists")
	ErrNotExist = errors.New("repository: does not exist")

	// storage errors
	ErrRetry = errors.New("repository: retry storage operation")

	// TODO: add metadata to this error? Like the column(s) causing the conflict?
	ErrConflict = errors.New("repository: conflict in storage operation")
)
