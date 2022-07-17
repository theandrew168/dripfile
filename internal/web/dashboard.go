package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/validator"
)

type dashboardForm struct {
	validator.Validator
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
