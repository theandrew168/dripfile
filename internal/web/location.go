package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/validator"
)

type locationCreateForm struct {
	validator.Validator `form:"-"`

	Endpoint        string `form:"Endpoint"`
	BucketName      string `form:"BucketName"`
	AccessKeyID     string `form:"AccessKeyID"`
	SecretAccessKey string `form:"SecretAccessKey"`
}

type locationDeleteForm struct {
	validator.Validator `form:"-"`

	LocationID string `form:"LocationID"`
}

type locationData struct {
	Locations []model.Location
	Location  model.Location
	Form      locationCreateForm
}

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	page := "app/location/list.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	project := session.Account.Project
	locations, err := app.store.Location.ReadAllByProject(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := locationData{
		Locations: locations,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
	page := "app/location/read.html"

	id := flow.Param(r.Context(), "id")
	location, err := app.store.Location.Read(id)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	data := locationData{
		Location: location,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	page := "app/location/create.html"
	data := locationData{}
	app.render(w, r, page, data)
}

func (app *Application) handleLocationCreateForm(w http.ResponseWriter, r *http.Request) {
	page := "app/location/create.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form locationCreateForm
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
		data := locationData{
			Form: form,
		}
		app.render(w, r, page, data)
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

		data := locationData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	// verify connection
	err = conn.Ping()
	if err != nil {
		form.SetError(err.Error())

		data := locationData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	jsonInfo, err := json.Marshal(info)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: encrypt connection info
	nonce, err := app.box.Nonce()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	encryptedInfo := app.box.Encrypt(nonce, jsonInfo)

	name := info.Endpoint + "/" + info.BucketName
	project := session.Account.Project
	location := model.NewLocation(model.KindS3, name, encryptedInfo, project)
	err = app.store.Location.Create(&location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("location create", map[string]string{
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"location_id": location.ID,
	})
	http.Redirect(w, r, "/location/"+location.ID, http.StatusSeeOther)
}

func (app *Application) handleLocationDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form locationDeleteForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// TODO: assert id belongs to session->account->project
	// TODO: assert account role is owner, admin, or editor
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

	app.logger.Info("location delete", map[string]string{
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"location_id": location.ID,
	})
	http.Redirect(w, r, "/location", http.StatusSeeOther)
}
