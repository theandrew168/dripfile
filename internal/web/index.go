package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/html/web"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	err := app.html.Web.Index(w, web.IndexParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
