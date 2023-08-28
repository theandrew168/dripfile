package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"
	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/location"
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

type CreateLocationRequest struct {
	Kind string `json:"kind"`
}

type CreateMemoryLocationRequest struct {
	Kind string `json:"kind"`
}

type CreateS3LocationRequest struct {
	Kind            string `json:"kind"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	locations, err := app.locationRepo.List()
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

	location, err := app.locationRepo.Read(id)
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
	body := readBody(w, r)

	// read the body info a buffer since we'll be decoding it (into JSON) mulitple times
	b, err := io.ReadAll(body)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	// just read the location "kind" for now
	var req CreateLocationRequest
	err = readJSON(bytes.NewReader(b), &req, false)
	if err != nil {
		app.badRequestResponse(w, r, err.Error())
		return
	}

	id, _ := uuid.NewRandom()
	var loc *location.Location

	// read more specific details based on the "kind"
	if req.Kind == location.KindMemory {
		loc, err = location.NewMemory(id.String())
		if err != nil {
			app.badRequestResponse(w, r, err.Error())
			return
		}
	} else if req.Kind == location.KindS3 {
		var req CreateS3LocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err.Error())
			return
		}

		loc, err = location.NewS3(id.String(), req.Endpoint, req.Bucket, req.AccessKeyID, req.SecretAccessKey)
		if err != nil {
			app.badRequestResponse(w, r, err.Error())
			return
		}
	} else {
		app.badRequestResponse(w, r, "empty or invalid location kind")
		return
	}

	// create the new location
	err = app.locationRepo.Create(loc)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	l := Location{
		ID:   loc.ID(),
		Kind: loc.Kind(),
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
