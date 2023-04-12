package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/lnquy/cron"

	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/view/web"
)

// https://gist.github.com/jpluimers/6510369
var shortcuts = map[string]string{
	"@yearly":   "0 0 1 1 *",
	"@annually": "0 0 1 1 *",
	"@monthly":  "0 0 1 * *",
	"@weekly":   "0 0 * * 0",
	"@daily":    "0 0 * * *",
	"@midnight": "0 0 * * *",
	"@hourly":   "0 * * * *",
}

func (app *Application) handleScheduleList(w http.ResponseWriter, r *http.Request) {
	schedules, err := app.store.Schedule.ReadAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.ScheduleListParams{
		Schedules: schedules,
	}
	err = app.view.Web.ScheduleList(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleScheduleRead(w http.ResponseWriter, r *http.Request) {
	id := flow.Param(r.Context(), "id")
	schedule, err := app.store.Schedule.Read(id)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	params := web.ScheduleReadParams{
		Schedule: schedule,
	}
	err = app.view.Web.ScheduleRead(w, params)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleScheduleCreate(w http.ResponseWriter, r *http.Request) {
	err := app.view.Web.ScheduleCreate(w, web.ScheduleCreateParams{})
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func (app *Application) handleScheduleCreateForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form web.ScheduleCreateForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Expr, "Expr")

	if !form.Valid() {
		// re-render with errors
		params := web.ScheduleCreateParams{
			Form: form,
		}
		err := app.view.Web.ScheduleCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	// check for shortcut exprs
	expr := form.Expr
	if shortcut, ok := shortcuts[expr]; ok {
		expr = shortcut
	}

	dtor, err := cron.NewDescriptor()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	name, err := dtor.ToDescription(expr, cron.Locale_en)
	if err != nil {
		form.SetFieldError("Expr", "Invalid cron expression")

		// re-render with errors
		params := web.ScheduleCreateParams{
			Form: form,
		}
		err := app.view.Web.ScheduleCreate(w, params)
		if err != nil {
			app.serverErrorResponse(w, r, err)
			return
		}

		return
	}

	schedule := model.NewSchedule(name, expr)
	err = app.store.Schedule.Create(&schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("schedule create", map[string]string{
		"account_id":  session.Account.ID,
		"schedule_id": schedule.ID,
	})
	http.Redirect(w, r, "/schedule/"+schedule.ID, http.StatusSeeOther)
}

func (app *Application) handleScheduleDeleteForm(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form web.ScheduleDeleteForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	// TODO: assert account role is admin or editor
	schedule, err := app.store.Schedule.Read(form.ScheduleID)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.store.Schedule.Delete(schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("schedule delete", map[string]string{
		"account_id":  session.Account.ID,
		"schedule_id": schedule.ID,
	})
	http.Redirect(w, r, "/schedule", http.StatusSeeOther)
}
