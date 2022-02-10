package connection

import (
	"io"
)

type Connection interface {
	List() ([]string, error)
	Read(path string) (io.Reader, error)
	Write(path string, r io.Reader) error
}
