package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
)

// TODO: if empty, big center button with basic info and an iamge
// TODO: if not empty, list view with button up in top right
func (app *Application) handleLocationReadMany(w http.ResponseWriter, r *http.Request) {
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
		Category  string
		Locations []core.Location
	}{
		Category:  "location",
		Locations: locations,
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/read_many.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationRead(w http.ResponseWriter, r *http.Request) {
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
		Category string
		Location core.Location
	}{
		Category: "location",
		Location: location,
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/read.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationCreate(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "location",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"location/create.page.html",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleLocationCreateForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: support other types of info
	endpoint := r.PostFormValue("endpoint")
	accessKeyID := r.PostFormValue("access-key-id")
	secretAccessKey := r.PostFormValue("secret-access-key")
	bucketName := r.PostFormValue("bucket-name")

	info := core.S3Info{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		SecretAccessKey: secretAccessKey,
		BucketName:      bucketName,
	}

	b, err := json.Marshal(info)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	location := core.NewLocation(core.KindS3, string(b), session.Account.Project)
	err = app.storage.Location.Create(&location)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/location", http.StatusSeeOther)
}
