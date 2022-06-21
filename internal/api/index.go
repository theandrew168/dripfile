package api

import (
	"html/template"
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFS(app.template, "index.html.tmpl")
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}

	err = ts.Execute(w, nil)
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}
}
