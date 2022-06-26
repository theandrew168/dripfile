package web

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"site.layout.html",
		"index.page.html",
	}

	app.render(w, r, files, nil)
}
