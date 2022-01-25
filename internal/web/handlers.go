package web

import (
	"errors"
	"net/http"

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
		"index.page.tmpl",
		"base.layout.tmpl",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleDashboard(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"dashboard.page.tmpl",
		"base.layout.tmpl",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleLogin(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"login.page.tmpl",
		"base.layout.tmpl",
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

	account, err := app.storage.Account.ReadByEmail(email)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			// TODO: handle not exists (invalid user or pass)
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

	// TODO: create new session linked to account
	// TODO: set session cookie

	app.logger.Info("login %s\n", email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"register.page.tmpl",
		"base.layout.tmpl",
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

	err = app.storage.Account.Create(&account)
	if err != nil {
		if errors.Is(err, core.ErrExist) {
			// TODO: handle exists
			app.serverErrorResponse(w, r, err)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("register %s\n", email)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
