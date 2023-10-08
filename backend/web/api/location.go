package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
)

func (app *Application) handleLocationCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleLocationCreate"))
	}
}

func (app *Application) handleLocationList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleLocationList"))
	}
}

func (app *Application) handleLocationRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleLocationRead: %s", id)
	}
}

func (app *Application) handleLocationDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleLocationDelete: %s", id)
	}
}
