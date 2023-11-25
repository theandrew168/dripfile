package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/validator"
)

type Itinerary struct {
	ID uuid.UUID `json:"id"`

	FromLocationID uuid.UUID `json:"fromLocationID"`
	ToLocationID   uuid.UUID `json:"toLocationID"`
	Pattern        string    `json:"pattern"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (app *Application) handleItineraryCreate() http.HandlerFunc {
	type request struct {
		FromLocationID string `json:"fromLocationID"`
		ToLocationID   string `json:"toLocationID"`
		Pattern        string `json:"pattern"`
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
		v.Check(req.FromLocationID != "", "fromLocationID", "must be provided")
		v.Check(req.ToLocationID != "", "toLocationID", "must be provided")
		v.Check(req.Pattern != "", "pattern", "must be provided")

		// check if provided IDs are valid UUIDs
		fromLocationID, err := uuid.Parse(req.FromLocationID)
		if err != nil {
			v.AddError("fromLocationID", "must be the ID of an existing location")
		}
		toLocationID, err := uuid.Parse(req.ToLocationID)
		if err != nil {
			v.AddError("toLocationID", "must be the ID of an existing location")
		}

		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		from, err := app.repo.Location.Read(fromLocationID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				v.AddError("fromLocationID", "must be the ID of an existing location")
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		to, err := app.repo.Location.Read(toLocationID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				v.AddError("toLocationID", "must be the ID of an existing location")
			default:
				app.serverErrorResponse(w, r, err)
				return
			}
		}

		// check if provided IDs correspond to existing entities
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		itinerary, err := domain.NewItinerary(from, to, req.Pattern)
		if err != nil {
			v.AddError("itinerary", err.Error())
		}

		// ensure new itinerary satisfies domain constraints
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.repo.Itinerary.Create(itinerary)
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
			ID: itinerary.ID(),

			FromLocationID: itinerary.FromLocationID(),
			ToLocationID:   itinerary.ToLocationID(),
			Pattern:        itinerary.Pattern(),
			CreatedAt:      itinerary.CreatedAt(),
			UpdatedAt:      itinerary.UpdatedAt(),
		}
		resp := response{
			Itinerary: apiItinerary,
		}

		header := make(http.Header)
		header.Set("Location", fmt.Sprintf("/api/v1/itinerary/%s", itinerary.ID()))

		err = writeJSON(w, http.StatusCreated, resp, header)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleItineraryList() http.HandlerFunc {
	type response struct {
		Itineraries []Itinerary `json:"itineraries"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		itineraries, err := app.repo.Itinerary.List()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// use make here to encode JSON as "[]" instead of "null" if empty
		apiItineraries := make([]Itinerary, 0)
		for _, itinerary := range itineraries {
			apiItinerary := Itinerary{
				ID: itinerary.ID(),

				FromLocationID: itinerary.FromLocationID(),
				ToLocationID:   itinerary.ToLocationID(),
				Pattern:        itinerary.Pattern(),
				CreatedAt:      itinerary.CreatedAt(),
				UpdatedAt:      itinerary.UpdatedAt(),
			}
			apiItineraries = append(apiItineraries, apiItinerary)
		}

		resp := response{
			Itineraries: apiItineraries,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleItineraryRead() http.HandlerFunc {
	type response struct {
		Itinerary Itinerary `json:"itinerary"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		itinerary, err := app.repo.Itinerary.Read(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		apiItinerary := Itinerary{
			ID: itinerary.ID(),

			FromLocationID: itinerary.FromLocationID(),
			ToLocationID:   itinerary.ToLocationID(),
			Pattern:        itinerary.Pattern(),
			CreatedAt:      itinerary.CreatedAt(),
			UpdatedAt:      itinerary.UpdatedAt(),
		}
		resp := response{
			Itinerary: apiItinerary,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleItineraryDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		itinerary, err := app.repo.Itinerary.Read(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		err = app.repo.Itinerary.Delete(itinerary)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			case errors.Is(err, domain.ErrLocationInUse):
				app.conflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
