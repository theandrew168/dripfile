package template_test

import (
	"bytes"
	"testing"
	"testing/fstest"

	"github.com/theandrew168/dripfile/internal/template"
	"github.com/theandrew168/dripfile/internal/test"
)

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
