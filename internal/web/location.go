package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/fileserver"
	"github.com/theandrew168/dripfile/internal/form"
)

func (app *Application) handleLocationList(w http.ResponseWriter, r *http.Request) {
	page := "location/list.page.html"

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

	data := struct {
		Locations []core.Location
	}{
		Locations: locations,
	}

	app.render(w, r, page, data)
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
	page := "location/read.page.html"

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

	data := struct {
		Location core.Location
	}{
		Location: location,
	}

	app.render(w, r, page, data)
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	page := "location/create.page.html"

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, page, data)
}

func (app *Application) handleLocationCreateForm(w http.ResponseWriter, r *http.Request) {
	page := "location/create.page.html"

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
		app.render(w, r, page, data)
		return
	}

	// TODO: support other types of info
	endpoint := f.Get("endpoint")
	bucketName := f.Get("bucket-name")
	accessKeyID := f.Get("access-key-id")
	secretAccessKey := f.Get("secret-access-key")

	info := fileserver.S3Info{
		Endpoint:        endpoint,
		BucketName:      bucketName,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
	}

	conn, err := fileserver.NewS3(info)
	if err != nil {
		data.Form.Errors.Add("general", err.Error())
		app.render(w, r, page, data)
		return
	}

	// verify connection
	err = conn.Ping()
	if err != nil {
		data.Form.Errors.Add("general", err.Error())
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
	location := core.NewLocation(core.KindS3, name, encryptedInfo, project)
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

	// TODO: assert id belongs to session->account->project
	// TODO: assert account role is owner, admin, or editor
	id := r.PostForm.Get("id")

	location, err := app.store.Location.Read(id)
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
