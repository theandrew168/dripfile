package repository

import "errors"

// TODO: add metadata to errors to make em more useful:
//   - what already exists
//   - what was missin
//   - what column(s) caused the conflict
var (
	// based the os package error names:
	// https://pkg.go.dev/os#pkg-variables
	ErrExist    = errors.New("repository: already exists")
	ErrNotExist = errors.New("repository: does not exist")

	// storage errors
	ErrRetry    = errors.New("repository: retry storage operation")
	ErrConflict = errors.New("repository: conflict in storage operation")
)
