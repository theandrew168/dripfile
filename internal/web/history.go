package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/model"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	page := "app/history/list.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	histories, err := app.store.History.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: check which xfer IDs are still valid
	// TODO: add map of valid IDs to tmpl data

	data := struct {
		Histories []model.History
	}{
		Histories: histories,
	}

	app.render(w, r, page, data)
}
