package memory

import (
	"bytes"
	"io"

	"github.com/theandrew168/dripfile/internal/location/fileserver"
)

type Info struct{}

func (info Info) Validate() error {
	return nil
}

type FileServer struct {
	info Info
	data map[string]*bytes.Buffer
}

func New(info Info) (*FileServer, error) {
	fs := FileServer{
		info: info,
		data: make(map[string]*bytes.Buffer),
	}

	return &fs, nil
}

func (fs *FileServer) Ping() error {
	return nil
}

func (fs *FileServer) Search(pattern string) ([]fileserver.FileInfo, error) {
	// TODO: implement this
	return []fileserver.FileInfo{}, nil
}

func (fs *FileServer) Read(file fileserver.FileInfo) (io.Reader, error) {
	// TODO: implement this
	return nil, nil
}

func (fs *FileServer) Write(file fileserver.FileInfo, r io.Reader) error {
	// TODO: implement this
	return nil
}
