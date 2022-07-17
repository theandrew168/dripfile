package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/validator"
)

type loginForm struct {
	validator.Validator `form:"-"`

	Email    string `form:"Email"`
	Password string `form:"Password"`
}

type loginData struct {
	Form loginForm
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/login.html"
	data := loginData{}
	app.render(w, r, page, data)
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/login.html"

	var form loginForm
	err := app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Email, "Email")
	form.CheckRequired(form.Password, "Password")

	if !form.Valid() {
		data := loginData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	account, err := app.store.Account.ReadByEmail(form.Email)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			form.SetError("Invalid email or password")

			data := loginData{
				Form: form,
			}
			app.render(w, r, page, data)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(form.Password))
	if err != nil {
		form.SetError("Invalid email or password")

		data := loginData{
			Form: form,
		}
		app.render(w, r, page, data)
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

	// create session model and store in the postgresql
	session := core.NewSession(sessionHash, expiry, account)
	err = app.store.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set permanent session cookie
	cookie := NewPermanentCookie(sessionIDCookieName, sessionID, expiry)
	http.SetCookie(w, &cookie)

	app.logger.Info("account login", map[string]string{
		"project_id": session.Account.Project.ID,
		"account_id": session.Account.ID,
	})
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
