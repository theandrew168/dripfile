package fileserver

import (
	"bytes"
	"errors"
	"io"
	"path/filepath"
	"sync"
)

var ErrNotFound = errors.New("fileserver: not found")

// ensure FileServer interface is satisfied
var _ FileServer = (*MemoryFileServer)(nil)

type MemoryInfo struct{}

type file struct {
	info FileInfo
	data *bytes.Buffer
}

type MemoryFileServer struct {
	sync.RWMutex
	info  MemoryInfo
	files map[string]file
}

func NewMemory(info MemoryInfo) (*MemoryFileServer, error) {
	fs := MemoryFileServer{
		info:  info,
		files: make(map[string]file),
	}

	return &fs, nil
}

func (fs *MemoryFileServer) Ping() error {
	return nil
}

func (fs *MemoryFileServer) Search(pattern string) ([]FileInfo, error) {
	fs.RLock()
	defer fs.RUnlock()

	var files []FileInfo
	for name, file := range fs.files {
		matched, _ := filepath.Match(pattern, name)
		if !matched {
			continue
		}

		files = append(files, file.info)
	}

	return files, nil
}

func (fs *MemoryFileServer) Read(name string) (io.Reader, error) {
	fs.RLock()
	defer fs.RUnlock()

	file, ok := fs.files[name]
	if !ok {
		return nil, ErrNotFound
	}

	return file.data, nil
}

func (fs *MemoryFileServer) Write(info FileInfo, r io.Reader) error {
	fs.Lock()
	defer fs.Unlock()

	buf, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	fs.files[info.Name] = file{
		info: info,
		data: bytes.NewBuffer(buf),
	}

	return nil
}
