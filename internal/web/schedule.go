package web

import (
	"net/http"
)

func (app *Application) handleScheduleReadMany(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "schedule",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/read_many.page.html",
	}

	app.render(w, r, files, data)
}
