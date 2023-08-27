package fileserver

import (
	"bytes"
	"io"
)

type MemoryInfo struct{}

func (info MemoryInfo) Validate() error {
	return nil
}

type memoryFileServer struct {
	info MemoryInfo
	data map[string]*bytes.Buffer
}

func NewMemory(info MemoryInfo) (FileServer, error) {
	fs := memoryFileServer{
		info: info,
		data: make(map[string]*bytes.Buffer),
	}

	return &fs, nil
}

func (fs *memoryFileServer) Ping() error {
	return nil
}

func (fs *memoryFileServer) Search(pattern string) ([]FileInfo, error) {
	// TODO: implement this
	return []FileInfo{}, nil
}

func (fs *memoryFileServer) Read(file FileInfo) (io.Reader, error) {
	// TODO: implement this
	return nil, nil
}

func (fs *memoryFileServer) Write(file FileInfo, r io.Reader) error {
	// TODO: implement this
	return nil
}
