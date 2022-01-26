package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/theandrew168/dripfile/internal/core"
)

// Redirects:
// 303 See Other         - for GETs after POSTs (like a login / register form)
// 302 Found             - all other temporary redirects
// 301 Moved Permanently - permanent redirects

// Route Handler Naming Ideas:
//
// basic page handlers:
// GET - handleIndex
// GET - handleDashboard
//
// basic page w/ form handlers:
// GET  - handleLogin
// POST - handleLoginForm
//
// CRUD handlers:
// C POST   - handleCreateFoo
// R GET    - handleReadFoo  (one / single / detail)
// R GET    - handleReadFoos (all / multiple / paginated)
// U PUT    - handleUpdateFoo
// D DELETE - handleDeleteFoo

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"index.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"register.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleRegisterForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	email := r.PostFormValue("email")
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	password = string(hash)
	account := core.NewAccount(email, username, password)

	// TODO: combine these storage ops into an atomic transaction somehow

	// create the account within the database
	err = app.storage.Account.Create(&account)
	if err != nil {
		if errors.Is(err, core.ErrExist) {
			// TODO: handle email exists
			app.serverErrorResponse(w, r, err)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	// create the new account's default project
	project := core.NewProject(account.Username)
	err = app.storage.Project.Create(&project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// link account <-> project with role "Owner"
	member := core.NewMember(core.RoleOwner, account, project)
	err = app.storage.Member.Create(&member)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// generate a fresh session ID
	sessionID, err := GenerateSessionID()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create session model and store in the database
	expiry := time.Now().AddDate(0, 0, 7)
	session := core.NewSession(sessionID, expiry, account)
	err = app.storage.Session.Create(&session)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// set cookie (just a session cookie after registration)
	cookie := NewSessionCookie(sessionIDCookieName, sessionID)
	http.SetCookie(w, &cookie)

	app.logger.Info("register %s\n", account.Email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"login.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleLoginForm(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	email := r.PostFormValue("email")
	password := r.PostFormValue("password")
	rememberMe := r.PostFormValue("remember-me")

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

	// create session model and store in the database
	expiry := time.Now().AddDate(0, 0, 7)
	session := core.NewSession(sessionID, expiry, account)
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

	app.logger.Info("login %s\n", account.Email)
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
	session, err := app.storage.Session.Read(sessionID.Value)
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

	app.logger.Info("logout %s\n", session.Account.Email)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	session, ok := r.Context().Value(contextKeySession).(core.Session)
	if !ok {
		err := fmt.Errorf("failed context value cast to core.Session")
		app.serverErrorResponse(w, r, err)
		return
	}
	app.logger.Info("%+v\n", session)

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"dashboard.page.html",
	}

	app.render(w, r, files, nil)
}
