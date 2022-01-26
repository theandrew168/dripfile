package web

import (
	"bytes"
	"html/template"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, html string) {
	files := []string{
		html,
		"base.layout.html",
	}

	// attempt to parse error template
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.Execute(&buf, nil)
	if err != nil {
		app.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// write the status and error page
	w.WriteHeader(status)
	w.Write(buf.Bytes())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusNotFound, "404.page.html")
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	app.errorResponse(w, r, http.StatusMethodNotAllowed, "405.page.html")
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err)
	app.errorResponse(w, r, http.StatusInternalServerError, "500.page.html")
}
