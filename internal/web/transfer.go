package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/form"
	"github.com/theandrew168/dripfile/internal/task"
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

	transfers, err := app.store.Transfer.ReadAllByProject(session.Account.Project)
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
	transfer, err := app.store.Transfer.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
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

	locations, err := app.store.Location.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.store.Schedule.ReadAllByProject(session.Account.Project)
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

	locations, err := app.store.Location.ReadAllByProject(session.Account.Project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.store.Schedule.ReadAllByProject(session.Account.Project)
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

	src, err := app.store.Location.Read(srcID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	dst, err := app.store.Location.Read(dstID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	schedule, err := app.store.Schedule.Read(scheduleID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	transfer := core.NewTransfer(pattern, src, dst, schedule, session.Account.Project)
	err = app.store.Transfer.Create(&transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("transfer create", map[string]string{
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"transfer_id": transfer.ID,
		"src_id":      transfer.Src.ID,
		"dst_id":      transfer.Dst.ID,
	})
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

	transfer, err := app.store.Transfer.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// submit this transfer to the task queue
	t, err := task.NewTransferTryTask(transfer.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	_, err = app.queue.Enqueue(t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("transfer run", map[string]string{
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"transfer_id": transfer.ID,
		"src_id":      transfer.Src.ID,
		"dst_id":      transfer.Dst.ID,
	})
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

	transfer, err := app.store.Transfer.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.store.Transfer.Delete(transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("transfer delete", map[string]string{
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"transfer_id": transfer.ID,
		"src_id":      transfer.Src.ID,
		"dst_id":      transfer.Dst.ID,
	})
	http.Redirect(w, r, "/transfer", http.StatusSeeOther)
}
