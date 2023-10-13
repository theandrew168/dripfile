package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"
	"github.com/theandrew168/dripfile/backend/model"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/validator"
)

type Itinerary struct {
	ID uuid.UUID `json:"id"`

	Pattern        string    `json:"pattern"`
	FromLocationID uuid.UUID `json:"from_location_id"`
	ToLocationID   uuid.UUID `json:"to_location_id"`
}

func (app *Application) handleItineraryCreate() http.HandlerFunc {
	type request struct {
		Pattern        string `json:"pattern"`
		FromLocationID string `json:"from_location_id"`
		ToLocationID   string `json:"to_location_id"`
	}

	type response struct {
		Itinerary Itinerary `json:"itinerary"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		body := readBody(w, r)

		var req request
		err := readJSON(body, &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// check if provided info passes basic validation
		v.Check(req.Pattern != "", "pattern", "must be provided")
		v.Check(req.FromLocationID != "", "from_location_id", "must be provided")
		v.Check(req.ToLocationID != "", "to_location_id", "must be provided")

		// check if provided IDs are valid UUIDs
		fromLocationID, err := uuid.Parse(req.FromLocationID)
		if err != nil {
			v.AddError("from_location_id", "must be a valid UUID")
		}
		toLocationID, err := uuid.Parse(req.ToLocationID)
		if err != nil {
			v.AddError("to_location_id", "must be a valid UUID")
		}

		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		itinerary := model.NewItinerary(req.Pattern, fromLocationID, toLocationID)

		err = app.srvc.Itinerary.Create(itinerary)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrConflict):
				app.conflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		apiItinerary := Itinerary{
			ID:             itinerary.ID,
			Pattern:        itinerary.Pattern,
			FromLocationID: itinerary.FromLocationID,
			ToLocationID:   itinerary.ToLocationID,
		}
		resp := response{
			Itinerary: apiItinerary,
		}

		header := make(http.Header)
		header.Set("Location", fmt.Sprintf("/api/v1/itinerary/%s", itinerary.ID))

		err = writeJSON(w, http.StatusCreated, resp, header)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
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
