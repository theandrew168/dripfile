package web

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	page := "index.page.html"
	app.render(w, r, page, nil)
}
