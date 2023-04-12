package web

import (
	"embed"
	"io/fs"
	"os"

	"github.com/theandrew168/dripfile/internal/html/template"
)

//go:embed template
var templateFS embed.FS

type HTML struct {
	reader *template.Reader
}

func New(reload bool) *HTML {
	var files fs.FS
	if reload {
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		files = os.DirFS("./internal/html/web/template/")
	} else {
		// else use the embedded template dir
		files, _ = fs.Sub(templateFS, "template")
	}

	html := HTML{
		reader: template.NewReader(files, reload),
	}
	return &html
}
