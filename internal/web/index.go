package web

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	page := "site/index.html"
	app.render(w, r, page, nil)
}
