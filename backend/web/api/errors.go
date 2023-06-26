package api

import (
	"net/http"
)

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}

	err := writeJSON(w, status, env)
	if err != nil {
		app.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	code := http.StatusBadRequest
	app.errorResponse(w, r, code, errors)
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	text := http.StatusText(code)
	app.errorResponse(w, r, code, text)
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	text := http.StatusText(code)
	app.errorResponse(w, r, code, text)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())

	code := http.StatusInternalServerError
	text := http.StatusText(code)
	app.errorResponse(w, r, code, text)
}
