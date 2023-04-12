package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/view/web"
)

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := app.view.Web.AuthLogin(w, web.AuthLoginParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	var form web.AuthLoginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Email, "Email")
	form.CheckRequired(form.Password, "Password")

	if !form.Valid() {
		// re-render with errors
		params := web.AuthLoginParams{
			Form: form,
		}
		err := app.view.Web.AuthLogin(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	account, err := app.store.Account.ReadByEmail(form.Email)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			form.SetError("Invalid email or password")

			// re-render with errors
			params := web.AuthLoginParams{
				Form: form,
			}
			err := app.view.Web.AuthLogin(w, params)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(form.Password))
	if err != nil {
		form.SetError("Invalid email or password")

		// re-render with errors
		params := web.AuthLoginParams{
			Form: form,
		}
		err := app.view.Web.AuthLogin(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

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

	// set permanent session cookie
	cookie := NewPermanentCookie(sessionIDCookieName, sessionID, expiry)
	http.SetCookie(w, &cookie)

	app.logger.Info("account login",
		slog.String("account_id", session.Account.ID),
	)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
