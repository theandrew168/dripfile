package api

import (
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error any `json:"error"`
}

func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	resp := ErrorResponse{
		Error: message,
	}

	err := writeJSON(w, status, resp, nil)
	if err != nil {
		app.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	code := http.StatusBadRequest
	app.errorResponse(w, r, code, err.Error())
}

func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusNotFound
	text := http.StatusText(code)
	app.errorResponse(w, r, code, strings.ToLower(text))
}

func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusMethodNotAllowed
	text := http.StatusText(code)
	app.errorResponse(w, r, code, strings.ToLower(text))
}

func (app *Application) conflictResponse(w http.ResponseWriter, r *http.Request) {
	code := http.StatusConflict
	text := http.StatusText(code)
	app.errorResponse(w, r, code, strings.ToLower(text))
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in our Validator type.
func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	code := http.StatusUnprocessableEntity
	app.errorResponse(w, r, code, errors)
}

func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error(err.Error())

	code := http.StatusInternalServerError
	text := http.StatusText(code)
	app.errorResponse(w, r, code, strings.ToLower(text))
}
