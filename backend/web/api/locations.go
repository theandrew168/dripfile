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

	err = writeJSON(w, 200, envelope{"locations": ls})
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

	err = writeJSON(w, 200, envelope{"location": l})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
