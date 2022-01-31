package web

import (
	"net/http"
)

func (app *Application) handleScheduleList(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "schedule",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/list.page.html",
	}

	app.render(w, r, files, data)
}
