package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/form"
	"github.com/theandrew168/dripfile/pkg/task"
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

	transfers, err := app.storage.Transfer.ReadAllByProject(session.Account.Project)
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

	locations, err := app.storage.Location.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.storage.Schedule.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Locations []core.Location
		Schedules []core.Schedule
		Form      *form.Form
	}{
		Locations: locations,
		Schedules: schedules,
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

	locations, err := app.storage.Location.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.storage.Schedule.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	f := form.New(r.PostForm)
	f.Required("pattern", "src-id", "dst-id")

	data := struct {
		Locations []core.Location
		Schedules []core.Schedule
		Form      *form.Form
	}{
		Locations: locations,
		Schedules: schedules,
		Form:      f,
	}

	if !f.Valid() {
		app.render(w, r, files, data)
		return
	}

	pattern := f.Get("pattern")
	srcID := f.Get("src-id")
	dstID := f.Get("dst-id")
	scheduleID := f.Get("schedule-id")

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

	schedule, err := app.storage.Schedule.Read(scheduleID)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	transfer := core.NewTransfer(pattern, src, dst, schedule, session.Account.Project)
	err = app.storage.Transfer.Create(&transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Printf("account %s transfer %s create\n", session.Account.Email, transfer.ID)
	http.Redirect(w, r, "/transfer/"+transfer.ID, http.StatusSeeOther)
}

func (app *Application) handleTransferRunForm(w http.ResponseWriter, r *http.Request) {
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

	// submit this transfer to the task queue
	t, err := task.DoTransfer(transfer.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.queue.Push(t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Printf("account %s transfer %s run\n", session.Account.Email, transfer.ID)
	http.Redirect(w, r, "/history", http.StatusSeeOther)
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

	app.infoLog.Printf("account %s transfer %s delete\n", session.Account.Email, transfer.ID)
	http.Redirect(w, r, "/transfer", http.StatusSeeOther)
}
