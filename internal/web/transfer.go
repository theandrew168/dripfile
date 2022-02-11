package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/core"
)

func (app *Application) handleTransferList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	transfers, err := app.storage.Transfer.ReadManyByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Category  string
		Transfers []core.Transfer
	}{
		Category:  "transfer",
		Transfers: transfers,
	}

	app.render(w, r, files, data)
}
