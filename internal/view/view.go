package view

import (
	"github.com/theandrew168/dripfile/internal/view/web"
)

type View struct {
	Web *web.View
}

func New(reload bool) *View {
	t := View{
		Web: web.New(reload),
	}
	return &t
}
