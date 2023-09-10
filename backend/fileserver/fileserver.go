package fileserver

import "io"

type FileInfo struct {
	Name string
	Size int
}

type FileServer interface {
	Ping() error
	Search(pattern string) ([]FileInfo, error)
	Read(name string) (io.Reader, error)
	Write(info FileInfo, r io.Reader) error
}
