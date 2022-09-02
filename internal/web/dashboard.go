package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/html/web"
)

// TODO: show some useful data on the dashboard
func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	err := app.html.Web.Dashboard(w, web.DashboardParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
