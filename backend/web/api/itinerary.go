package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
)

func (app *Application) handleItineraryCreate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleItineraryCreate"))
	}
}

func (app *Application) handleItineraryList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("TODO: handleItineraryList"))
	}
}

func (app *Application) handleItineraryRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleItineraryRead: %s", id)
	}
}

func (app *Application) handleItineraryDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleItineraryDelete: %s", id)
	}
}

func (app *Application) handleItineraryRun() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := flow.Param(r.Context(), "id")
		fmt.Fprintf(w, "TODO: handleItineraryRun: %s", id)
	}
}
