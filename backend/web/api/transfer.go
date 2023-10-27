package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/validator"
)

type Transfer struct {
	ID uuid.UUID `json:"id"`

	ItineraryID uuid.UUID             `json:"itineraryID"`
	Status      domain.TransferStatus `json:"status"`
	Progress    int                   `json:"progress"`
}

func (app *Application) handleTransferCreate() http.HandlerFunc {
	type request struct {
		ItineraryID string `json:"itineraryID"`
	}

	type response struct {
		Transfer Transfer `json:"transfer"`
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
		v.Check(req.ItineraryID != "", "itineraryID", "must be provided")

		// check if provided IDs are valid UUIDs
		itineraryID, err := uuid.Parse(req.ItineraryID)
		if err != nil {
			v.AddError("itineraryID", "must be the ID of an existing location")
		}

		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		itinerary, err := app.repo.Itinerary.Read(itineraryID)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				v.AddError("itineraryID", "must be the ID of an existing location")
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

		transfer, err := domain.NewTransfer(itinerary)
		if err != nil {
			v.AddError("transfer", err.Error())
		}

		// ensure new transfer satisfies domain constraints
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		err = app.repo.Transfer.Create(transfer)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrConflict):
				app.conflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		// TODO: kick off transfer in background

		apiTransfer := Transfer{
			ID: transfer.ID(),

			ItineraryID: transfer.ItineraryID(),
			Status:      transfer.Status(),
			Progress:    transfer.Progress(),
		}
		resp := response{
			Transfer: apiTransfer,
		}

		header := make(http.Header)
		header.Set("Location", fmt.Sprintf("/api/v1/transfers/%s", transfer.ID()))

		err = writeJSON(w, http.StatusCreated, resp, header)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleTransferList() http.HandlerFunc {
	type response struct {
		Transfers []Transfer `json:"transfers"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		transfers, err := app.repo.Transfer.List()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// use make here to encode JSON as "[]" instead of "null" if empty
		apiTransfers := make([]Transfer, 0)
		for _, transfer := range transfers {
			apiTransfer := Transfer{
				ID: transfer.ID(),

				ItineraryID: transfer.ItineraryID(),
				Status:      transfer.Status(),
				Progress:    transfer.Progress(),
			}
			apiTransfers = append(apiTransfers, apiTransfer)
		}

		resp := response{
			Transfers: apiTransfers,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleTransferRead() http.HandlerFunc {
	type response struct {
		Transfer Transfer `json:"transfer"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		transfer, err := app.repo.Transfer.Read(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		apiTransfer := Transfer{
			ID: transfer.ID(),

			ItineraryID: transfer.ItineraryID(),
			Status:      transfer.Status(),
			Progress:    transfer.Progress(),
		}
		resp := response{
			Transfer: apiTransfer,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}
