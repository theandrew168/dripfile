package web

import (
	"bytes"
	"net/http"

	"github.com/theandrew168/dripfile/internal/html/errors"
)

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request) {
	// render template to a temp buffer
	var b bytes.Buffer
	err := app.html.Errors.Error400(&b, errors.Error400Params{})
	if err != nil {
		app.logger.Error(err, nil)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// write the status and error page
	w.WriteHeader(http.StatusBadRequest)
	w.Write(b.Bytes())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	// render template to a temp buffer
	var b bytes.Buffer
	err := app.html.Errors.Error404(&b, errors.Error404Params{})
	if err != nil {
		app.logger.Error(err, nil)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// write the status and error page
	w.WriteHeader(http.StatusNotFound)
	w.Write(b.Bytes())
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	// render template to a temp buffer
	var b bytes.Buffer
	err := app.html.Errors.Error405(&b, errors.Error405Params{})
	if err != nil {
		app.logger.Error(err, nil)

		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// write the status and error page
	w.WriteHeader(http.StatusMethodNotAllowed)
	w.Write(b.Bytes())
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	// log details of the error locally but the user sees a generic 500
	app.logger.Error(err, nil)

	// render template to a temp buffer
	var b bytes.Buffer
	err = app.html.Errors.Error500(&b, errors.Error500Params{})
	if err != nil {
		code := http.StatusInternalServerError
		http.Error(w, http.StatusText(code), code)
		return
	}

	// write the status and error page
	w.WriteHeader(http.StatusInternalServerError)
	w.Write(b.Bytes())
}
