package fileserver_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/test"
)

func TestPing(t *testing.T) {
	fs, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	err = fs.Ping()
	test.AssertNilError(t, err)
}

func TestSearch(t *testing.T) {
	fs, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	contents := "testing"
	info := fileserver.FileInfo{
		Name: "foo.txt",
		Size: int64(len(contents)),
	}

	err = fs.Write(info, bytes.NewBufferString(contents))
	test.AssertNilError(t, err)

	infos, err := fs.Search("*.txt")
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(infos), 1)
	test.AssertSliceContains(t, infos, info)
}

func TestRead(t *testing.T) {
	fs, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	contents := "testing"
	info := fileserver.FileInfo{
		Name: "foo.txt",
		Size: int64(len(contents)),
	}

	err = fs.Write(info, bytes.NewBufferString(contents))
	test.AssertNilError(t, err)

	r, err := fs.Read("foo.txt")
	test.AssertNilError(t, err)

	buf, err := io.ReadAll(r)
	test.AssertNilError(t, err)

	test.AssertEqual(t, string(buf), contents)
}

func TestWrite(t *testing.T) {
	fs, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	info := fileserver.FileInfo{
		Name: "foo.txt",
		Size: 42,
	}
	data := bytes.NewBufferString("testing")

	err = fs.Write(info, data)
	test.AssertNilError(t, err)
}
