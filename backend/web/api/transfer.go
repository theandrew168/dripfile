package api

import (
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"
)

func (app *Application) handleTransferCreate() http.HandlerFunc {
	type request struct {
		itineraryID uuid.UUID `json:"itinerary_id"`
	}
	type response struct {
		transferID uuid.UUID `json:"transfer_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		// lookup itinerary
		// create transfer
		// kick off transfer in background
		// return transfer ID
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
