package web

import (
	"net/http"
)

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
	files := []string{
		"login.page.tmpl",
		"base.layout.tmpl",
	}

	app.render(w, r, files, nil)
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

	app.logger.Info(email)
	app.logger.Info(username)
	app.logger.Info(password)

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
