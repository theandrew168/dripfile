package fileserver

import "io"

type FileInfo struct {
	Name string
	Size int64
}

type FileServer interface {
	Ping() error
	Search(pattern string) ([]FileInfo, error)
	Read(path string) (io.Reader, error)
	Write(info FileInfo, r io.Reader) error
}
