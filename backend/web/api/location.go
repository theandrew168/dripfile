package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/google/uuid"

	"github.com/theandrew168/dripfile/backend/domain"
	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/validator"
)

type Location struct {
	ID uuid.UUID `json:"id"`

	Kind string `json:"kind"`
}

func (app *Application) handleLocationCreate() http.HandlerFunc {
	type request struct {
		Kind string `json:"kind"`
	}
	type requestMemory struct {
		Kind string `json:"kind"`
	}
	type requestS3 struct {
		Kind            string `json:"kind"`
		Endpoint        string `json:"endpoint"`
		Bucket          string `json:"bucket"`
		AccessKeyID     string `json:"accessKeyID"`
		SecretAccessKey string `json:"secretAccessKey"`
	}

	type response struct {
		Location Location `json:"location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		v := validator.New()
		body := readBody(w, r)

		// read the body info a buffer since we'll be decoding it (into JSON) mulitple times
		b, err := io.ReadAll(body)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		// just read the location "kind" for now
		var req request
		err = readJSON(bytes.NewReader(b), &req, false)
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		kind := domain.LocationKind(req.Kind)
		v.Check(
			validator.PermittedValue(kind, domain.LocationKindMemory, domain.LocationKindS3),
			"kind",
			"must be one of: memory, s3",
		)
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		// read more specific details based on the "kind"
		var location *domain.Location
		if kind == domain.LocationKindMemory {
			var req requestMemory
			err = readJSON(bytes.NewReader(b), &req, true)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			if !v.Valid() {
				app.failedValidationResponse(w, r, v.Errors)
				return
			}

			location, err = domain.NewMemoryLocation()
			if err != nil {
				v.AddError("location", err.Error())
			}
		} else if kind == domain.LocationKindS3 {
			var req requestS3
			err = readJSON(bytes.NewReader(b), &req, true)
			if err != nil {
				app.badRequestResponse(w, r, err)
				return
			}

			v.Check(req.Endpoint != "", "endpoint", "must be provided")
			v.Check(req.Bucket != "", "bucket", "must be provided")
			v.Check(req.AccessKeyID != "", "accessKeyID", "must be provided")
			v.Check(req.SecretAccessKey != "", "secretAccessKey", "must be provided")
			if !v.Valid() {
				app.failedValidationResponse(w, r, v.Errors)
				return
			}

			location, err = domain.NewS3Location(
				req.Endpoint,
				req.Bucket,
				req.AccessKeyID,
				req.SecretAccessKey,
			)
			if err != nil {
				v.AddError("location", err.Error())
			}
		}

		// ensure new location satisfies domain constraints
		if !v.Valid() {
			app.failedValidationResponse(w, r, v.Errors)
			return
		}

		// create the new location
		err = app.repo.Location.Create(location)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrConflict):
				app.conflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		apiLocation := Location{
			ID:   location.ID(),
			Kind: string(location.Kind()),
		}
		resp := response{
			Location: apiLocation,
		}

		header := make(http.Header)
		header.Set("Location", fmt.Sprintf("/api/v1/locations/%s", location.ID()))

		err = writeJSON(w, http.StatusCreated, resp, header)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleLocationList() http.HandlerFunc {
	type response struct {
		Locations []Location `json:"locations"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		locations, err := app.repo.Location.List()
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		// use make here to encode JSON as "[]" instead of "null" if empty
		apiLocations := make([]Location, 0)
		for _, location := range locations {
			apiLocation := Location{
				ID:   location.ID(),
				Kind: string(location.Kind()),
			}
			apiLocations = append(apiLocations, apiLocation)
		}

		resp := response{
			Locations: apiLocations,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleLocationRead() http.HandlerFunc {
	type response struct {
		Location Location `json:"location"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		location, err := app.repo.Location.Read(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		apiLocation := Location{
			ID:   location.ID(),
			Kind: string(location.Kind()),
		}
		resp := response{
			Location: apiLocation,
		}

		err = writeJSON(w, http.StatusOK, resp, nil)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}
	}
}

func (app *Application) handleLocationDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(flow.Param(r.Context(), "id"))
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		location, err := app.repo.Location.Read(id)
		if err != nil {
			switch {
			case errors.Is(err, repository.ErrNotExist):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}

			return
		}

		err = app.repo.Location.Delete(location)
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
