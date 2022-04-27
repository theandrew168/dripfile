package web

import (
	"net/http"
)

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := struct{}{}

	files := []string{
		"app.layout.html",
		"dashboard.page.html",
	}

	app.render(w, r, files, data)
}
