package web

import (
	"fmt"
	"net/http"

	"github.com/theandrew168/dripfile/internal/core"
)

func (app *Application) handleHistoryList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"history/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	histories, err := app.storage.History.ReadManyByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, h := range histories {
		fmt.Printf("%+v\n", h)
	}

	data := struct {
		Histories []core.History
	}{
		Histories: histories,
	}

	app.render(w, r, files, data)
}
