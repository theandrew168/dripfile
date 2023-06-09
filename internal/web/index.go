package web

import (
	"net/http"
)

func (app *Application) handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello API"))
}
