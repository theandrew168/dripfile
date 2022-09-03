package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/view/web"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	history, err := app.store.History.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: check which xfer IDs are still valid
	// TODO: add map of valid IDs to tmpl data

	params := web.HistoryListParams{
		History: history,
	}
	err = app.view.Web.HistoryList(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
