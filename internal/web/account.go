package web

import (
	"net/http"

	"github.com/theandrew168/dripfile/internal/view/web"
)

func (app *Application) handleAccountRead(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.AccountReadParams{
		Account: session.Account,
	}
	err = app.view.Web.AccountRead(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleAccountDeleteForm(w http.ResponseWriter, r *http.Request) {
	var form web.AccountDeleteForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if form.AccountID != session.Account.ID {
		app.badRequestResponse(w, r)
		return
	}

	err = app.store.Account.Delete(session.Account)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
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
