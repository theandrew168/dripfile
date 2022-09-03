package view

import (
	"github.com/theandrew168/dripfile/internal/view/api"
	"github.com/theandrew168/dripfile/internal/view/web"
)

type Template struct {
	API *api.Template
	Web *web.Template
}

func New(reload bool) *Template {
	t := Template{
		API: api.New(reload),
		Web: web.New(reload),
	}
	return &t
}
