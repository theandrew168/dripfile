package core

import (
	"io"
)

// abstraction over backend locations (S3, SFTP, etc)
type Connection interface {
	List() ([]string, error)
	Read(path string) (io.Reader, error)
	Write(path string, r io.Reader) error
}
