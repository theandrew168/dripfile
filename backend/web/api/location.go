package api

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/theandrew168/dripfile/backend/database"
)

type Location struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

type SingleLocationResponse struct {
	Location Location `json:"location"`
}

type MultipleLocationResponse struct {
	Locations []Location `json:"locations"`
}

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	locations, err := app.locationStorage.List()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var ls []Location
	for _, location := range locations {
		l := Location{
			ID:   location.ID(),
			Kind: location.Kind(),
		}
		ls = append(ls, l)
	}

	resp := MultipleLocationResponse{
		Locations: ls,
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	location, err := app.locationStorage.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	l := Location{
		ID:   location.ID(),
		Kind: location.Kind(),
	}

	resp := SingleLocationResponse{
		Location: l,
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	location, err := app.locationStorage.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	l := Location{
		ID:   location.ID(),
		Kind: location.Kind(),
	}

	resp := SingleLocationResponse{
		Location: l,
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
