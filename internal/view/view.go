package view

import (
	"github.com/theandrew168/dripfile/internal/view/api"
	"github.com/theandrew168/dripfile/internal/view/web"
)

type View struct {
	API *api.View
	Web *web.View
}

func New(reload bool) *View {
	t := View{
		API: api.New(reload),
		Web: web.New(reload),
	}
	return &t
}
