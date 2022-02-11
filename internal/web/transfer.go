package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/form"
)

func (app *Application) handleTransferList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	transfers, err := app.storage.Transfer.ReadManyByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Transfers []core.Transfer
	}{
		Transfers: transfers,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleTransferRead(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/read.page.html",
	}

	id := flow.Param(r.Context(), "id")
	transfer, err := app.storage.Transfer.Read(id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Transfer core.Transfer
	}{
		Transfer: transfer,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleTransferCreate(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/create.page.html",
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
		Form      *form.Form
	}{
		Locations: locations,
		Form:      form.New(nil),
	}

	app.render(w, r, files, data)
}

func (app *Application) handleTransferCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/create.page.html",
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

	f := form.New(r.PostForm)
	f.Required("pattern", "src-id", "dst-id")

	data := struct {
		Locations []core.Location
		Form      *form.Form
	}{
		Locations: locations,
		Form:      f,
	}

	if !f.Valid() {
		app.render(w, r, files, data)
		return
	}

	pattern := r.PostForm.Get("pattern")
	srcID := r.PostForm.Get("src-id")
	dstID := r.PostForm.Get("dst-id")

	src, err := app.storage.Location.Read(srcID)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	dst, err := app.storage.Location.Read(dstID)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	transfer := core.NewTransfer(pattern, src, dst, session.Account.Project)
	err = app.storage.Transfer.Create(&transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("account %s transfer %s create\n", session.Account.Email, transfer.ID)
	http.Redirect(w, r, "/transfer/"+transfer.ID, http.StatusSeeOther)
}

func (app *Application) handleTransferDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: assert id belongs to session->account->project
	// TODO: assert account role is owner, admin, or editor
	id := r.PostForm.Get("id")

	transfer, err := app.storage.Transfer.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.storage.Transfer.Delete(transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("account %s transfer %s delete\n", session.Account.Email, transfer.ID)
	http.Redirect(w, r, "/transfer", http.StatusSeeOther)
}
