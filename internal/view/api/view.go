package api

import (
	"embed"
	"io/fs"
	"os"

	"github.com/theandrew168/dripfile/internal/view/template"
)

//go:embed template
var templateFS embed.FS

type View struct {
	r *template.Reader
}

func New(reload bool) *View {
	var files fs.FS
	if reload {
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		files = os.DirFS("./internal/view/api/template/")
	} else {
		// else use the embedded template dir
		files, _ = fs.Sub(templateFS, "template")
	}

	v := View{
		r: template.NewReader(files, reload),
	}
	return &v
}
