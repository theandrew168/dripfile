package web

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/form"
)

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"auth/register.page.html",
	}

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, files, data)
}

func (app *Application) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"auth/register.page.html",
	}

	f := form.New(r.PostForm)
	f.Required("email", "username", "password")

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
		app.render(w, r, files, data)
		return
	}

	email := f.Get("email")
	username := f.Get("username")
	password := f.Get("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: combine these storage ops into an atomic transaction somehow

	// create project for the new account
	project := core.NewProject()
	err = app.storage.Project.Create(&project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create the new account
	account := core.NewAccount(email, username, string(hash), core.RoleOwner, project)
	err = app.storage.Account.Create(&account)
	if err != nil {
		if errors.Is(err, core.ErrExist) {
			f.Errors.Add("email", "An account with this email already exists")
			app.render(w, r, files, data)
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
	session := core.NewSession(sessionHash, expiry, account)
	err = app.storage.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set cookie (just a session cookie after registration)
	cookie := NewSessionCookie(sessionIDCookieName, sessionID)
	http.SetCookie(w, &cookie)

	app.logger.Info("account %s create\n", account.Email)
	app.logger.Info("account %s login\n", account.Email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"auth/login.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")
	rememberMe := r.PostForm.Get("remember-me")

	account, err := app.storage.Account.ReadByEmail(email)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			// TODO: handle email not exists (invalid user or pass)
			app.serverErrorResponse(w, r, err)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		app.logger.Info("invalid creds!\n")
		// TODO: handle invalid creds (invalid user or pass)
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
	session := core.NewSession(sessionHash, expiry, account)
	err = app.storage.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set cookie (session / permanent based on "Remember me")
	if rememberMe != "" {
		cookie := NewPermanentCookie(sessionIDCookieName, sessionID, expiry)
		http.SetCookie(w, &cookie)
	} else {
		cookie := NewSessionCookie(sessionIDCookieName, sessionID)
		http.SetCookie(w, &cookie)
	}

	app.logger.Info("account %s login\n", account.Email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleLogoutForm(w http.ResponseWriter, r *http.Request) {
	// check for session cookie
	sessionID, err := r.Cookie(sessionIDCookieName)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// check for session in database
	sessionHash := fmt.Sprintf("%x", sha256.Sum256([]byte(sessionID.Value)))
	session, err := app.storage.Session.Read(sessionHash)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// delete session from database
	err = app.storage.Session.Delete(session)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// expire the existing session cookie
	cookie := NewExpiredCookie(sessionIDCookieName)
	http.SetCookie(w, &cookie)

	app.logger.Info("account %s logout\n", session.Account.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
