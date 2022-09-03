package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
)

type HistoryListParams struct {
	History []model.History
}

func (t *Template) HistoryList(w io.Writer, p HistoryListParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"history/list.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
