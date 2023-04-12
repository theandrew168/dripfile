package html

import (
	"github.com/theandrew168/dripfile/internal/html/web"
)

// TODO: site, app, docs, api, etc? split based on first part of path?
type View struct {
	Web *web.View
}

func New(reload bool) *View {
	t := View{
		Web: web.New(reload),
	}
	return &t
}
