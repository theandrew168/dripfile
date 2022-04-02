package web

import (
	"bytes"
	"html/template"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, code int, files []string) {
	// parse template files
	ts, err := template.ParseFS(app.templates, files...)
	if err != nil {
		app.errorLog.Println(err)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err = ts.ExecuteTemplate(&buf, "base", nil)
	if err != nil {
		app.errorLog.Println(err)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// write the status and error page
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(code)
	w.Write(buf.Bytes())
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"error/400.page.html",
	}
	app.errorResponse(w, r, http.StatusBadRequest, files)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"error/404.page.html",
	}
	app.errorResponse(w, r, http.StatusNotFound, files)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"error/405.page.html",
	}
	app.errorResponse(w, r, http.StatusMethodNotAllowed, files)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log details of the error locally but the user sees a generic 500
	app.errorLog.Println(err)

	files := []string{
		"base.layout.html",
		"error/500.page.html",
	}
	app.errorResponse(w, r, http.StatusInternalServerError, files)
}
