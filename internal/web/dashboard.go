package web

import (
	"net/http"
)

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	page := "dashboard.page.html"
	app.render(w, r, page, nil)
}
