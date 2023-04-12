package html

import (
	"github.com/theandrew168/dripfile/internal/html/web"
)

// TODO: site, app, docs, api, etc? split based on first part of path?
type HTML struct {
	Web *web.HTML
}

func New(reload bool) *HTML {
	html := HTML{
		Web: web.New(reload),
	}
	return &html
}
