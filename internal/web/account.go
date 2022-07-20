package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/model"
)

type accountData struct {
	Account model.Account
}

func (app *Application) handleAccountRead(w http.ResponseWriter, r *http.Request) {
	page := "app/account/read.html"
	data := accountData{}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data.Account = session.Account
	app.render(w, r, page, data)
}

func (app *Application) handleAccountDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: use form AccountID or just session.Account?
	err = app.store.Account.Delete(session.Account)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// check number of accounts linked to project
	count, err := app.store.Account.CountByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: keep this behavior?
	// delete project if no more associated accounts
	if count == 0 {
		err = app.store.Project.Delete(session.Account.Project)
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
