package api

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/transfer"
	"github.com/theandrew168/dripfile/backend/validator"
)

type Transfer struct {
	ID             string `json:"id"`
	Pattern        string `json:"pattern"`
	FromLocationID string `json:"from_location_id"`
	ToLocationID   string `json:"to_location_id"`
}

type SingleTransferResponse struct {
	Transfer Transfer `json:"transfer"`
}

type MultipleTransferResponse struct {
	Transfers []Transfer `json:"transfers"`
}

type CreateTransferRequest struct {
	Pattern        string `json:"pattern"`
	FromLocationID string `json:"from_location_id"`
	ToLocationID   string `json:"to_location_id"`
}

func (app *Application) handleTransferCreate(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	body := readBody(w, r)

	var req CreateTransferRequest
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
	fromLocationId, err := uuid.Parse(req.FromLocationID)
	if err != nil {
		v.AddError("from_location_id", "must be a valid UUID")
	}
	toLocationId, err := uuid.Parse(req.ToLocationID)
	if err != nil {
		v.AddError("to_location_id", "must be a valid UUID")
	}

	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fromLocation, err := app.locationRepo.Read(fromLocationId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			v.AddError("from_location_id", "must be the ID of an existing location")
		default:
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	toLocation, err := app.locationRepo.Read(toLocationId)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			v.AddError("to_location_id", "must be the ID of an existing location")
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

	id := uuid.Must(uuid.NewRandom())
	t, err := transfer.New(id, req.Pattern, fromLocation, toLocation)
	if err != nil {
		switch {
		case errors.Is(err, transfer.ErrInvalidPattern):
			v.AddError("pattern", "must be a valid pattern")
		case errors.Is(err, transfer.ErrSameLocation):
			v.AddError("to_location_id", "must be different than from_location_id")
		default:
			v.AddError("transfer", err.Error())
		}
	}

	// check if new transfer is valid
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// create the new transfer
	err = app.transferRepo.Create(t)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrConflict):
			app.conflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	resp := SingleTransferResponse{
		Transfer: Transfer{
			ID:             t.ID().String(),
			Pattern:        t.Pattern(),
			FromLocationID: t.FromLocationID().String(),
			ToLocationID:   t.ToLocationID().String(),
		},
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferList(w http.ResponseWriter, r *http.Request) {
	transfers, err := app.transferRepo.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var ts []Transfer
	for _, transfer := range transfers {
		t := Transfer{
			ID:             transfer.ID().String(),
			Pattern:        transfer.Pattern(),
			FromLocationID: transfer.FromLocationID().String(),
			ToLocationID:   transfer.ToLocationID().String(),
		}
		ts = append(ts, t)
	}

	resp := MultipleTransferResponse{
		Transfers: ts,
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferRead(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	transfer, err := app.transferRepo.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	t := Transfer{
		ID:             transfer.ID().String(),
		Pattern:        transfer.Pattern(),
		FromLocationID: transfer.FromLocationID().String(),
		ToLocationID:   transfer.ToLocationID().String(),
	}

	resp := SingleTransferResponse{
		Transfer: t,
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.transferRepo.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
