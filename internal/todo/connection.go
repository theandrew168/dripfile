package core

import (
	"io"
)

type Connection interface {
	Verify() bool
	List() ([]string, error)
	Read(path string) (io.Reader, error)
	Write(path string, r io.Reader) error
}
