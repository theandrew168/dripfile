package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/lnquy/cron"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/postgresql"
	"github.com/theandrew168/dripfile/internal/validator"
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

type scheduleForm struct {
	validator.Validator `form:"-"`

	Expr string `form:"Expr"`
}

type scheduleData struct {
	Schedules []core.Schedule
	Schedule  core.Schedule
	Form      scheduleForm
}

func (app *Application) handleScheduleList(w http.ResponseWriter, r *http.Request) {
	page := "app/schedule/list.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	project := session.Account.Project
	schedules, err := app.store.Schedule.ReadAllByProject(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := scheduleData{
		Schedules: schedules,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleScheduleRead(w http.ResponseWriter, r *http.Request) {
	page := "app/schedule/read.html"

	id := flow.Param(r.Context(), "id")
	schedule, err := app.store.Schedule.Read(id)
	if err != nil {
		if errors.Is(err, postgresql.ErrNotExist) {
			app.notFoundResponse(w, r)
			return
		}

		app.serverErrorResponse(w, r, err)
		return
	}

	data := scheduleData{
		Schedule: schedule,
	}
	app.render(w, r, page, data)
}

func (app *Application) handleScheduleCreate(w http.ResponseWriter, r *http.Request) {
	page := "app/schedule/create.html"
	data := scheduleData{}
	app.render(w, r, page, data)
}

func (app *Application) handleScheduleCreateForm(w http.ResponseWriter, r *http.Request) {
	page := "app/schedule/create.html"

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	dtor, err := cron.NewDescriptor()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var form scheduleForm
	err = app.decodePostForm(r, &form)
	if err != nil {
		app.badRequestResponse(w, r)
		return
	}

	form.CheckRequired(form.Expr, "Expr")

	if !form.Valid() {
		data := scheduleData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	// check for shortcut exprs
	expr := form.Expr
	if shortcut, ok := shortcuts[expr]; ok {
		expr = shortcut
	}

	name, err := dtor.ToDescription(expr, cron.Locale_en)
	if err != nil {
		form.SetFieldError("Expr", "Invalid cron expression")
		data := scheduleData{
			Form: form,
		}
		app.render(w, r, page, data)
		return
	}

	project := session.Account.Project
	schedule := core.NewSchedule(name, expr, project)
	err = app.store.Schedule.Create(&schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.logger.Info("schedule create", map[string]string{
		"project_id":  session.Account.Project.ID,
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

	// TODO: assert id belongs to session->account->project
	// TODO: assert account role is owner, admin, or editor
	scheduleID := r.PostForm.Get("ScheduleID")
	schedule, err := app.store.Schedule.Read(scheduleID)
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
		"project_id":  session.Account.Project.ID,
		"account_id":  session.Account.ID,
		"schedule_id": schedule.ID,
	})
	http.Redirect(w, r, "/schedule", http.StatusSeeOther)
}
