package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/lnquy/cron"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/database"
	"github.com/theandrew168/dripfile/internal/form"
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
	page := "schedule/list.page.html"

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

	data := struct {
		Schedules []core.Schedule
	}{
		Schedules: schedules,
	}

	app.render(w, r, page, data)
}

func (app *Application) handleScheduleRead(w http.ResponseWriter, r *http.Request) {
	page := "schedule/read.page.html"

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

	data := struct {
		Schedule core.Schedule
	}{
		Schedule: schedule,
	}

	app.render(w, r, page, data)
}

func (app *Application) handleScheduleCreate(w http.ResponseWriter, r *http.Request) {
	page := "schedule/create.page.html"

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, page, data)
}

func (app *Application) handleScheduleCreateForm(w http.ResponseWriter, r *http.Request) {
	page := "schedule/create.page.html"

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

	f := form.New(r.PostForm)
	f.Required("expr")

	expr := f.Get("expr")
	if shortcut, ok := shortcuts[expr]; ok {
		expr = shortcut
	}

	name, err := dtor.ToDescription(expr, cron.Locale_en)
	if err != nil {
		f.Errors.Add("expr", "invalid cron expression")
	}

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
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
	id := r.PostForm.Get("id")

	schedule, err := app.store.Schedule.Read(id)
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
