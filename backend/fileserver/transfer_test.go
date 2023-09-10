package fileserver_test

import (
	"bytes"
	"io"
	"testing"

	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/test"
)

// TODO: Run tests for each FileServer impl

func TestTransfer(t *testing.T) {
	random := test.NewRandom()

	// connect to a source location
	from, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	// create two files on the "from" location
	name := "foo.txt"
	size := 20
	contents := random.String(size)

	err = from.Write(
		fileserver.FileInfo{Name: name, Size: size},
		bytes.NewBufferString(contents),
	)
	test.AssertNilError(t, err)

	err = from.Write(
		fileserver.FileInfo{Name: "foo.png", Size: size},
		bytes.NewBuffer(random.Bytes(size)),
	)
	test.AssertNilError(t, err)

	// connect to a destination location
	to, err := fileserver.NewMemory(fileserver.MemoryInfo{})
	test.AssertNilError(t, err)

	// run the transfer
	totalBytes, err := fileserver.Transfer("*.txt", from, to)
	test.AssertNilError(t, err)
	test.AssertEqual(t, totalBytes, size)

	// ensure only one file was copied
	files, err := to.Search("*")
	test.AssertNilError(t, err)
	test.AssertEqual(t, len(files), 1)

	file := files[0]
	test.AssertEqual(t, file.Name, name)
	test.AssertEqual(t, file.Size, size)

	// read the file and verify contents
	r, err := to.Read(name)
	test.AssertNilError(t, err)

	buf, err := io.ReadAll(r)
	test.AssertNilError(t, err)
	test.AssertEqual(t, string(buf), contents)
}
