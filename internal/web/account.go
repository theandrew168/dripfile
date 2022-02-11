package web

import (
	"net/http"
)

func (app *Application) handleAccountRead(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"account/read.page.html",
	}

	app.render(w, r, files, nil)
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
	project := session.Account.Project
	count, err := app.storage.Account.CountByProject(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// delete project if no more associated accounts
	if count == 0 {
		err = app.storage.Project.Delete(project)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// expire the existing session cookie
	cookie := NewExpiredCookie(sessionIDCookieName)
	http.SetCookie(w, &cookie)

	app.logger.Info("account %s delete\n", session.Account.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
