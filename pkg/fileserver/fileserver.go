package fileserver

import (
	"io"
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
	Rename(src, dst FileInfo) error
	Delete(file FileInfo) error
}
