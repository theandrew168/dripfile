package web

// TODO: renderTemplate helper?

import (
	"html/template"
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.tmpl",
	}

	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
