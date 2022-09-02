package api

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/html/api"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := app.html.API.Index(w, api.IndexParams{})
	if err != nil {
		// TODO: JSON response w/ error message
		http.Error(w, "Internal server error", 500)
		return
	}
}
