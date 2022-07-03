package web

import (
	"bytes"
	"fmt"
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, code int, page string) {
	ts, ok := app.tm[page]
	if !ok {
		err := fmt.Errorf("web: template does not exist: %s", page)
		app.logger.Error(err, nil)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// render template to a temp buffer
	var buf bytes.Buffer
	err := ts.ExecuteTemplate(&buf, "base", nil)
	if err != nil {
		app.logger.Error(err, nil)

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
	page := "error/400.html"
	app.errorResponse(w, r, http.StatusBadRequest, page)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	page := "error/404.html"
	app.errorResponse(w, r, http.StatusNotFound, page)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	page := "error/405.html"
	app.errorResponse(w, r, http.StatusMethodNotAllowed, page)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log details of the error locally but the user sees a generic 500
	app.logger.Error(err, nil)

	page := "error/500.html"
	app.errorResponse(w, r, http.StatusInternalServerError, page)
}
