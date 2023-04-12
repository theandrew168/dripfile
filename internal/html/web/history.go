package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
)

type HistoryListParams struct {
	History []model.History
}

func (html *HTML) HistoryList(w io.Writer, p HistoryListParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"history/list.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}
