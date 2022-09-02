package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/html/site"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := app.html.Site.Index(w, site.IndexParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
