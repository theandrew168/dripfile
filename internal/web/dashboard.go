package web

import (
	"net/http"
)

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "dashboard",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"dashboard.page.html",
	}

	app.render(w, r, files, data)
}
