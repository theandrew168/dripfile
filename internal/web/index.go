package web

import (
	"net/http"
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
