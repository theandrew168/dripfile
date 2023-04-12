package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/html/web"
	"github.com/theandrew168/dripfile/internal/model"
)

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	err := app.html.Web.AuthRegister(w, web.AuthRegisterParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	var form web.AuthRegisterForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Email, "Email")
	form.CheckRequired(form.Password, "Password")

	if !form.Valid() {
		// re-render with errors
		params := web.AuthRegisterParams{
			Form: form,
		}
		err := app.html.Web.AuthRegister(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	account, err := app.srvc.CreateAccount(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, database.ErrExist) {
			form.SetFieldError("Email", "An account with this email already exists")

			// re-render with errors
			params := web.AuthRegisterParams{
				Form: form,
			}
			err := app.html.Web.AuthRegister(w, params)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	// generate a fresh session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID)))
	expiry := time.Now().AddDate(0, 0, 7)

	// create session model and store in the database
	session := model.NewSession(sessionHash, expiry, account)
	err = app.store.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set cookie (just a session cookie after registration)
	cookie := NewSessionCookie(sessionIDCookieName, sessionID)
	http.SetCookie(w, &cookie)

	app.logger.Info("account create",
		slog.String("account_id", session.Account.ID),
	)
	app.logger.Info("account login",
		slog.String("account_id", session.Account.ID),
	)

	// redirect to dashboard
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
