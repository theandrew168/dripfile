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

// abstraction over backend file servers (S3, SFTP, etc)
type FileServer interface {
	List() ([]string, error)
	Read(path string) (io.Reader, error)
	Write(path string, r io.Reader) error
}

type FTPInfo struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type FTPSInfo struct {
	Endpoint string `json:"endpoint"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type SFTPInfo struct {
	Endpoint   string `json:"endpoint"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
}
