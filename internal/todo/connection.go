package core

import (
	"io"
)

type Connection interface {
	Verify() bool
	List() ([]string, error)
	Read(name string) (io.Reader, error)
	Write(name string, r io.Reader) error
}
