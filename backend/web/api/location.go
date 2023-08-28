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

func (req *CreateLocationRequest) Validate(v *validator.Validator) {
	v.Check(validator.PermittedValue(req.Kind, location.KindMemory, location.KindS3), "kind", "must be one of: memory, s3")
}

type CreateMemoryLocationRequest struct {
	Kind string `json:"kind"`
}

func (req *CreateMemoryLocationRequest) Validate(v *validator.Validator) {}

type CreateS3LocationRequest struct {
	Kind            string `json:"kind"`
	Endpoint        string `json:"endpoint"`
	Bucket          string `json:"bucket"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

func (req *CreateS3LocationRequest) Validate(v *validator.Validator) {
	v.Check(req.Endpoint != "", "endpoint", "must be provided")
	v.Check(req.Bucket != "", "bucket", "must be provided")
	v.Check(req.AccessKeyID != "", "access_key_id", "must be provided")
	v.Check(req.SecretAccessKey != "", "secret_access_key", "must be provided")
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

	req.Validate(v)
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	id, _ := uuid.NewRandom()
	var loc *location.Location

	// read more specific details based on the "kind"
	if req.Kind == location.KindMemory {
		var req CreateMemoryLocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		req.Validate(v)
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		loc, err = location.NewMemory(id.String())
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	} else if req.Kind == location.KindS3 {
		var req CreateS3LocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		req.Validate(v)
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		loc, err = location.NewS3(id.String(), req.Endpoint, req.Bucket, req.AccessKeyID, req.SecretAccessKey)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
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

func (app *Application) handleLocationUpdate(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	loc, err := app.locationRepo.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
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

	req.Validate(v)
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

		req.Validate(v)
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		info := fileserver.MemoryInfo{}
		err := loc.SetMemory(info)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	} else if req.Kind == location.KindS3 {
		var req CreateS3LocationRequest
		err = readJSON(bytes.NewReader(b), &req, true)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		req.Validate(v)
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
		err := loc.SetS3(info)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}
	}

	// update the existing location
	err = app.locationRepo.Update(loc)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

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

func (app *Application) handleLocationDelete(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")

	err := app.locationRepo.Delete(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}
}
