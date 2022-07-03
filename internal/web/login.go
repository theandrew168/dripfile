package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/form"
)

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/login.html"

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, page, data)
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	page := "site/auth/login.html"

	f := form.New(r.PostForm)
	f.Required("email", "password")

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
		app.render(w, r, page, data)
		return
	}

	email := f.Get("email")
	password := f.Get("password")

	account, err := app.store.Account.ReadByEmail(email)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			f.Errors.Add("email", "Invalid email")
			app.render(w, r, page, data)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		f.Errors.Add("password", "Invalid password")
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

	// create session model and store in the database
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
