package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/html/web"
	"github.com/theandrew168/dripfile/internal/model"
)

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	locations, err := app.store.Location.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.LocationListParams{
		Locations: locations,
	}
	err = app.view.Web.LocationList(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")
	location, err := app.store.Location.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.LocationReadParams{
		Location: location,
	}
	err = app.view.Web.LocationRead(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	err := app.view.Web.LocationCreate(w, web.LocationCreateParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleLocationCreateForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form web.LocationCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Endpoint, "Endpoint")
	form.CheckRequired(form.BucketName, "BucketName")
	form.CheckRequired(form.AccessKeyID, "AccessKeyID")
	form.CheckRequired(form.SecretAccessKey, "SecretAccessKey")

	if !form.Valid() {
		// re-render with errors
		params := web.LocationCreateParams{
			Form: form,
		}
		err := app.view.Web.LocationCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	info := fileserver.S3Info{
		Endpoint:        form.Endpoint,
		BucketName:      form.BucketName,
		AccessKeyID:     form.AccessKeyID,
		SecretAccessKey: form.SecretAccessKey,
	}

	conn, err := fileserver.NewS3(info)
	if err != nil {
		form.SetError(err.Error())

		// re-render with errors
		params := web.LocationCreateParams{
			Form: form,
		}
		err := app.view.Web.LocationCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	// verify connection
	err = conn.Ping()
	if err != nil {
		form.SetError(err.Error())

		// re-render with errors
		params := web.LocationCreateParams{
			Form: form,
		}
		err := app.view.Web.LocationCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// encrypt connection info
	nonce, err := app.box.Nonce()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	encryptedInfo := app.box.Encrypt(nonce, jsonInfo)

	// store location w/ encrypted info
	name := info.Endpoint + "/" + info.BucketName
	location := model.NewLocation(model.KindS3, name, encryptedInfo)
	err = app.store.Location.Create(&location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("location create",
		slog.String("account_id", session.Account.ID),
		slog.String("location_id", location.ID),
	)
	http.Redirect(w, r, "/location/"+location.ID, http.StatusSeeOther)
}

func (app *Application) handleLocationDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form web.LocationDeleteForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// TODO: assert account role is admin or editor
	location, err := app.store.Location.Read(form.LocationID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.store.Location.Delete(location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("location delete",
		slog.String("account_id", session.Account.ID),
		slog.String("location_id", location.ID),
	)
	http.Redirect(w, r, "/location", http.StatusSeeOther)
}
