package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/connection"
	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/form"
)

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	locations, err := app.storage.Location.ReadManyByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Locations []core.Location
	}{
		Locations: locations,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/read.page.html",
	}

	id := flow.Param(r.Context(), "id")
	location, err := app.storage.Location.Read(id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Location core.Location
	}{
		Location: location,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/create.page.html",
	}

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/create.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	f := form.New(r.PostForm)

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
		app.render(w, r, files, data)
		return
	}

	// TODO: support other types of info
	endpoint := r.PostForm.Get("endpoint")
	bucketName := r.PostForm.Get("bucket-name")
	accessKeyID := r.PostForm.Get("access-key-id")
	secretAccessKey := r.PostForm.Get("secret-access-key")

	info := connection.S3Info{
		Endpoint:        endpoint,
		BucketName:      bucketName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}

	conn, err := connection.NewS3(info)
	if err != nil {
		data.Form.Errors.Add("general", err.Error())
		app.render(w, r, files, data)
		return
	}

	// verify connection
	_, err = conn.List()
	if err != nil {
		data.Form.Errors.Add("general", err.Error())
		app.render(w, r, files, data)
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

	location := core.NewLocation(core.KindS3, encryptedInfo, session.Account.Project)
	err = app.storage.Location.Create(&location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("account %s location %s create\n", session.Account.Email, location.ID)
	http.Redirect(w, r, "/location/"+location.ID, http.StatusSeeOther)
}

func (app *Application) handleLocationDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: assert id belongs to session->account->project
	// TODO: assert account role is owner, admin, or editor
	id := r.PostForm.Get("id")

	location, err := app.storage.Location.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.storage.Location.Delete(location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("account %s location %s delete\n", session.Account.Email, location.ID)
	http.Redirect(w, r, "/location", http.StatusSeeOther)
}
