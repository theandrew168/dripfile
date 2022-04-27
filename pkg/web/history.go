package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/pkg/core"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"app.layout.html",
		"history/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	histories, err := app.storage.History.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: check which xfer IDs are still valid
	// TODO: add map of valid IDs to tmpl data

	data := struct {
		Histories []core.History
	}{
		Histories: histories,
	}

	app.render(w, r, files, data)
}
