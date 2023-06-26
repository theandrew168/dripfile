package api

import "net/http"

type Location struct {
	ID   string `json:"id"`
	Kind string `json:"kind"`
}

func (app *Application) handleListLocations(w http.ResponseWriter, r *http.Request) {
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
