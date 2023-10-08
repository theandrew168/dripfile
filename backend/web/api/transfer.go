package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
)

func (app *Application) handleTransferCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleTransferCreate"))
	}
}

func (app *Application) handleTransferList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleTransferList"))
	}
}

func (app *Application) handleTransferRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleTransferRead: %s", id)
	}
}
