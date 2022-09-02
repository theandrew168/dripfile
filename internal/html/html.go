package html

import (
	"github.com/theandrew168/dripfile/internal/html/site"
)

type Template struct {
	Site *site.Template
	// ...
}

func New() *Template {
	t := Template{
		Site: site.New(),
	}
	return &t
}
