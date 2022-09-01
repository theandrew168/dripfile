package web

import (
	"net/http"
)

// helper for rendering templates and handling potential errors
func (app *Application) render(w http.ResponseWriter, r *http.Request, page string, data any) {
	t, err := app.tmpl.Get(page)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// always execute the base template
	w.Header().Set("Content-Type", "text/html")
	err = t.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
