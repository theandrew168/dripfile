package web

import (
	"html/template"
	"net/http"
)

func (app *Application) render(w http.ResponseWriter, r *http.Request, files []string, data interface{}) {
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"index.page.tmpl",
		"base.layout.tmpl",
	}

	app.render(w, r, files, nil)
}
