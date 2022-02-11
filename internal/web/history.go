package web

import (
	"net/http"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"history/list.page.html",
	}

	app.render(w, r, files, data)
}
