package fileserver

import (
	"bytes"
	"io"
)

// ensure FileServer interface is satisfied
var _ FileServer = (*MemoryFileServer)(nil)

type MemoryInfo struct{}

func (info MemoryInfo) Validate() error {
	return nil
}

type MemoryFileServer struct {
	info MemoryInfo
	data map[string]*bytes.Buffer
}

func NewMemory(info MemoryInfo) (*MemoryFileServer, error) {
	fs := MemoryFileServer{
		info: info,
		data: make(map[string]*bytes.Buffer),
	}

	return &fs, nil
}

func (fs *MemoryFileServer) Ping() error {
	return nil
}

func (fs *MemoryFileServer) Search(pattern string) ([]FileInfo, error) {
	// TODO: implement this
	return []FileInfo{}, nil
}

func (fs *MemoryFileServer) Read(file FileInfo) (io.Reader, error) {
	// TODO: implement this
	return nil, nil
}

func (fs *MemoryFileServer) Write(file FileInfo, r io.Reader) error {
	// TODO: implement this
	return nil
}
