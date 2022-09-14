package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/view/web"
)

func (app *Application) handleTransferList(w http.ResponseWriter, r *http.Request) {
	transfers, err := app.store.Transfer.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.TransferListParams{
		Transfers: transfers,
	}
	err = app.view.Web.TransferList(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferRead(w http.ResponseWriter, r *http.Request) {
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

	params := web.TransferReadParams{
		Transfer: transfer,
	}
	err = app.view.Web.TransferRead(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferCreate(w http.ResponseWriter, r *http.Request) {
	locations, err := app.store.Location.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.store.Schedule.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.TransferCreateParams{
		Locations: locations,
		Schedules: schedules,
	}
	err = app.view.Web.TransferCreate(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleTransferCreateForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// reload locations and schedules in case form needs to be rerendered
	locations, err := app.store.Location.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	schedules, err := app.store.Schedule.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form web.TransferCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Pattern, "Pattern")

	if !form.Valid() {
		// re-render with errors
		params := web.TransferCreateParams{
			Form: form,

			Locations: locations,
			Schedules: schedules,
		}
		err := app.view.Web.TransferCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

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

	transfer := model.NewTransfer(form.Pattern, src, dst, schedule)
	err = app.store.Transfer.Create(&transfer)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("transfer create", map[string]string{
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

	var form web.TransferRunForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// TODO: assert account role is admin or editor
	transfer, err := app.store.Transfer.Read(form.TransferID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	// submit this transfer to the task queue
	t := task.NewTransferTryTask(transfer.ID)
	err = app.queue.Submit(t)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("transfer run", map[string]string{
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

	var form web.TransferDeleteForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// TODO: assert account role is admin or editor
	transfer, err := app.store.Transfer.Read(form.TransferID)
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
		"account_id":  session.Account.ID,
		"transfer_id": transfer.ID,
		"src_id":      transfer.Src.ID,
		"dst_id":      transfer.Dst.ID,
	})
	http.Redirect(w, r, "/transfer", http.StatusSeeOther)
}
