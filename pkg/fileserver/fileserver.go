package fileserver

import (
	"errors"
	"io"
)

var (
	ErrInvalidEndpoint    = errors.New("core: invalid endpoint")
	ErrInvalidCredentials = errors.New("core: invalid credentials")
	ErrInvalidBucket      = errors.New("core: invalid bucket")
)

type FileInfo struct {
	Name string
	Size int64
}

type FileServer interface {
	Ping() error
	Search(pattern string) ([]FileInfo, error)
	Read(file FileInfo) (io.Reader, error)
	Write(file FileInfo, r io.Reader) error
}
