package template_test

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/theandrew168/dripfile/internal/html/template"
	"github.com/theandrew168/dripfile/internal/test"
)

func TestReaderEmpty(t *testing.T) {
	t.Parallel()

	reader := template.NewReader(nil, false)

	tmpl := reader.Read()
	if tmpl != nil {
		t.Fatalf("got: %v; want: nil", tmpl)
	}
}

func TestReaderBasic(t *testing.T) {
	t.Parallel()

	files := fstest.MapFS{
		"index.html": {
			Data: []byte("<html><body>Hello World!</body></html>"),
		},
	}
	reader := template.NewReader(files, false)

	var buf bytes.Buffer
	tmpl := reader.Read("index.html")
	err := tmpl.Execute(&buf, nil)
	test.AssertNilError(t, err)

	test.AssertStringContains(t, buf.String(), "Hello World!")
}
