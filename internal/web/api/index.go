package api

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := app.index.Execute(w, nil)
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}
}
