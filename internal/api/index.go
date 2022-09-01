package api

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	page := "api/index.html"

	t, err := app.tmpl.Get(page)
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = t.ExecuteTemplate(w, "base", nil)
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}
}
