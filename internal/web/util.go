package web

import (
	"html/template"
	"net/http"
)

// helper for rendering templates and handling potential errors
func (app *Application) render(w http.ResponseWriter, r *http.Request, files []string, data interface{}) {
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// always execute the base template
	w.Header().Set("Content-Type", "text/html")
	err = ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
