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
)

type loginForm struct {
	Form
	Email    string
	Password string
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
	data := loginData{}

	form := loginForm{
		Email:    r.PostForm.Get("email"),
		Password: r.PostForm.Get("password"),
	}

	form.CheckField(NotBlank(form.Email), "email", "This field cannot be blank")
	form.CheckField(NotBlank(form.Password), "password", "This field cannot be blank")

	if !form.Valid() {
		data.Form = form
		app.render(w, r, page, data)
		return
	}

	account, err := app.store.Account.ReadByEmail(form.Email)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			form.AddError("general", "Invalid email or password")
			data.Form = form
			app.render(w, r, page, data)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(form.Password))
	if err != nil {
		form.AddError("general", "Invalid email or password")
		data.Form = form
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
