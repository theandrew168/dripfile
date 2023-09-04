package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/database"
	"github.com/theandrew168/dripfile/backend/fileserver"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/validator"
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

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	v := validator.New()
	body := readBody(w, r)

	// read the body info a buffer since we'll be decoding it (into JSON) mulitple times
	b, err := io.ReadAll(body)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// just read the location "kind" for now
	var req CreateLocationRequest
	err = readJSON(bytes.NewReader(b), &req, false)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v.Check(validator.PermittedValue(req.Kind, location.KindMemory, location.KindS3), "kind", "must be one of: memory, s3")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	id := uuid.Must(uuid.NewRandom())
	var l *location.Location

	// read more specific details based on the "kind"
	if req.Kind == location.KindMemory {
		var req CreateMemoryLocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// TODO: no checks for this yet
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		l, err = location.NewMemory(id)
	} else if req.Kind == location.KindS3 {
		var req CreateS3LocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		v.Check(req.Endpoint != "", "endpoint", "must be provided")
		v.Check(req.Bucket != "", "bucket", "must be provided")
		v.Check(req.AccessKeyID != "", "access_key_id", "must be provided")
		v.Check(req.SecretAccessKey != "", "secret_access_key", "must be provided")
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		l, err = location.NewS3(id, req.Endpoint, req.Bucket, req.AccessKeyID, req.SecretAccessKey)
	}
	if err != nil {
		switch {
		case errors.Is(err, location.ErrInvalidKind):
			v.AddError("kind", "must be one of: memory, s3")
		default:
			v.AddError("location", err.Error())
		}
	}

	// check if new location is valid
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// create the new location
	err = app.locationRepo.Create(l)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrConflict):
			app.conflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	resp := SingleLocationResponse{
		Location: Location{
			ID:   l.ID().String(),
			Kind: l.Kind(),
		},
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
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
			ID:   location.ID().String(),
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
	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	location, err := app.locationRepo.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	l := Location{
		ID:   location.ID().String(),
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

func (app *Application) handleLocationUpdate(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	l, err := app.locationRepo.Read(id)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrNotExist):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	v := validator.New()
	body := readBody(w, r)

	// read the body info a buffer since we'll be decoding it (into JSON) mulitple times
	b, err := io.ReadAll(body)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// just read the location "kind" for now
	var req CreateLocationRequest
	err = readJSON(bytes.NewReader(b), &req, false)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v.Check(validator.PermittedValue(req.Kind, location.KindMemory, location.KindS3), "kind", "must be one of: memory, s3")
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// read more specific details based on the "kind"
	if req.Kind == location.KindMemory {
		var req CreateMemoryLocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// TODO: no checks for this yet
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		info := fileserver.MemoryInfo{}
		err = l.SetMemory(info)
	} else if req.Kind == location.KindS3 {
		var req CreateS3LocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		v.Check(req.Endpoint != "", "endpoint", "must be provided")
		v.Check(req.Bucket != "", "bucket", "must be provided")
		v.Check(req.AccessKeyID != "", "access_key_id", "must be provided")
		v.Check(req.SecretAccessKey != "", "secret_access_key", "must be provided")
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		info := fileserver.S3Info{
			Endpoint:        req.Endpoint,
			Bucket:          req.Bucket,
			AccessKeyID:     req.AccessKeyID,
			SecretAccessKey: req.SecretAccessKey,
		}
		err = l.SetS3(info)
	}
	if err != nil {
		switch {
		case errors.Is(err, location.ErrInvalidKind):
			v.AddError("kind", "must be one of: memory, s3")
		default:
			v.AddError("location", err.Error())
		}
	}

	// check if new location is valid
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	// update the existing location
	err = app.locationRepo.Update(l)
	if err != nil {
		switch {
		case errors.Is(err, database.ErrConflict):
			app.conflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}

		return
	}

	resp := SingleLocationResponse{
		Location: Location{
			ID:   l.ID().String(),
			Kind: l.Kind(),
		},
	}
	err = writeJSON(w, 200, resp)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

// TODO: Domain Service: warn if this location is being used by any transfers?
func (app *Application) handleLocationDelete(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(flow.Param(r.Context(), "id"))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.locationRepo.Delete(id)
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
