package site

import (
	"embed"
	"io/fs"
	"os"

	"github.com/theandrew168/dripfile/internal/html/template"
)

//go:embed template
var templateFS embed.FS

type Template struct {
	template.Template
}

func New(reload bool) *Template {
	var files fs.FS
	if reload {
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		files = os.DirFS("./internal/html/site/template/")
	} else {
		// else use the embedded template dir
		files, _ = fs.Sub(templateFS, "template")
	}

	t := Template{
		template.New(files, reload),
	}
	return &t
}
