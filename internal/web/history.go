package web

import (
	"net/http"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "history",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"history/list.page.html",
	}

	app.render(w, r, files, data)
}
