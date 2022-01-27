package web

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
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
// C POST   - handleCreateFoo[Form]
// R GET    - handleReadFoo[s]
// U PUT    - handleUpdateFoo[Form]
// D DELETE - handleDeleteFoo[Form]

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"site.layout.html",
		"index.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleRegister(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"auth/register.page.html",
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
		"auth/login.page.html",
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

	data := struct {
		Category string
	}{
		Category: "dashboard",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"dashboard.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleReadTransfers(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "transfer",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/read_all.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleReadLocations(w http.ResponseWriter, r *http.Request) {
	// TODO: if empty, big center button with basic info and an iamge
	// TODO: if not empty, list view with button up in top right

	locations, err := app.storage.Location.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Category  string
		Locations []core.Location
	}{
		Category:  "location",
		Locations: locations,
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/read_all.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleCreateLocation(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "location",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/create.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleReadLocation(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	location, err := app.storage.Location.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	data := struct {
		Category string
		Location core.Location
	}{
		Category: "location",
		Location: location,
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/read.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleReadSchedules(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "schedule",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/read_all.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleReadHistory(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "history",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"history/read_all.page.html",
	}

	app.render(w, r, files, data)
}
