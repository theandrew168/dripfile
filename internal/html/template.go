package html

import (
	"github.com/theandrew168/dripfile/internal/html/api"
	"github.com/theandrew168/dripfile/internal/html/app"
	"github.com/theandrew168/dripfile/internal/html/site"
)

type Template struct {
	App  *app.Template
	API  *api.Template
	Site *site.Template
}

func New(reload bool) *Template {
	t := Template{
		App:  app.New(reload),
		API:  api.New(reload),
		Site: site.New(reload),
	}
	return &t
}
