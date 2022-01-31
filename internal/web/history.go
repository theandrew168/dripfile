package web

import (
	"net/http"
)

func (app *Application) handleHistoryReadMany(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "history",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"history/read_many.page.html",
	}

	app.render(w, r, files, data)
}
