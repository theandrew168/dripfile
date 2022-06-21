package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/core"
)

func (app *Application) handleAccountRead(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"account/read.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Account core.Account
	}{
		Account: session.Account,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleAccountDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.storage.Account.Delete(session.Account)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// check number of accounts linked to project
	count, err := app.storage.Account.CountByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// delete project if no more associated accounts
	if count == 0 {
		err = app.storage.Project.Delete(session.Account.Project)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// expire the existing session cookie
	cookie := NewExpiredCookie(sessionIDCookieName)
	http.SetCookie(w, &cookie)

	app.logger.Info("account delete", map[string]string{
		"project_id": session.Account.Project.ID,
		"account_id": session.Account.ID,
	})
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
