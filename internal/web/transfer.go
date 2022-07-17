package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/validator"
)

type transferForm struct {
	validator.Validator
	Pattern    string
	SrcID      string
	DstID      string
	ScheduleID string
}

type transferData struct {
	Locations []core.Location
	Schedules []core.Schedule
	Transfers []core.Transfer
	Transfer  core.Transfer
	Form      transferForm
}

func (app *Application) handleTransferList(w http.ResponseWriter, r *http.Request) {
	page := "app/transfer/list.html"

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

	data := transferData{
		Transfers: transfers,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleTransferRead(w http.ResponseWriter, r *http.Request) {
	page := "app/transfer/read.html"

	id := flow.Param(r.Context(), "id")
	transfer, err := app.store.Transfer.Read(id)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	data := transferData{
		Transfer: transfer,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleTransferCreate(w http.ResponseWriter, r *http.Request) {
	page := "app/transfer/create.html"

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

	data := transferData{
		Locations: locations,
		Schedules: schedules,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleTransferCreateForm(w http.ResponseWriter, r *http.Request) {
	page := "app/transfer/create.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// reload locations and schedules in case form needs to be rerendered
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

	form := transferForm{
		Pattern:    r.PostForm.Get("Pattern"),
		SrcID:      r.PostForm.Get("SrcID"),
		DstID:      r.PostForm.Get("DstID"),
		ScheduleID: r.PostForm.Get("ScheduleID"),
	}

	form.CheckRequired(form.Pattern, "Pattern")

	if !form.Valid() {
		data := transferData{
			Locations: locations,
			Schedules: schedules,
			Form:      form,
		}
		app.render(w, r, page, data)
		return
	}

	src, err := app.store.Location.Read(form.SrcID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	dst, err := app.store.Location.Read(form.DstID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	schedule, err := app.store.Schedule.Read(form.ScheduleID)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.badRequestResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	project := session.Account.Project
	transfer := core.NewTransfer(form.Pattern, src, dst, schedule, project)
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
	transferID := r.PostForm.Get("TransferID")
	transfer, err := app.store.Transfer.Read(transferID)
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

	_, err = app.asynqClient.Enqueue(t)
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
	transferID := r.PostForm.Get("TransferID")
	transfer, err := app.store.Transfer.Read(transferID)
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
