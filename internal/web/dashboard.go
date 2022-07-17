package web

import (
	"net/http"
)

type dashboardForm struct {
	Form
	Search string
}

type dashboardData struct {
	Form dashboardForm
}

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	page := "app/dashboard.html"
	data := dashboardData{}
	app.render(w, r, page, data)
}
