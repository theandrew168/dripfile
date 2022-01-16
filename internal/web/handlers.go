package web

import (
	"html/template"
	"net/http"
)

/*
Route Handler Naming Ideas:

static-ish, read-only handlers:
GET - handleIndex
GET - handleDashboard

CRUD handlers:
C POST   - createFoo
R GET    - readFoo  (one / single / detail)
R GET    - readFoos (all / multiple / paginated)
U PUT    - updateFoo
D DELETE - deleteFoo
*/

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

func (app *Application) handleDoLogin(w http.ResponseWriter, r *http.Request) {
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

func (app *Application) handleDoRegister(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"register.page.tmpl",
		"base.layout.tmpl",
	}

	app.render(w, r, files, nil)
}

// helper for rendering templates and handling potential errors
func (app *Application) render(w http.ResponseWriter, r *http.Request, files []string, data interface{}) {
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	err = ts.Execute(w, data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
