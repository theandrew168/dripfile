package web

import (
	"errors"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/lnquy/cron"

	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/form"
)

func (app *Application) handleScheduleList(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/list.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	project := session.Account.Project
	schedules, err := app.storage.Schedule.ReadManyByProject(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		Schedules []core.Schedule
	}{
		Schedules: schedules,
	}

	app.render(w, r, files, data)
}

func (app *Application) handleScheduleRead(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/read.page.html",
	}

	id := flow.Param(r.Context(), "id")
	schedule, err := app.storage.Schedule.Read(id)
	if err != nil {
		if errors.Is(err, core.ErrNotExist) {
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

	app.render(w, r, files, data)
}

func (app *Application) handleScheduleCreate(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/create.page.html",
	}

	data := struct {
		Form *form.Form
	}{
		Form: form.New(nil),
	}

	app.render(w, r, files, data)
}

func (app *Application) handleScheduleCreateForm(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"schedule/create.page.html",
	}

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

	name, err := dtor.ToDescription(f.Get("expr"), cron.Locale_en)
	if err != nil {
		f.Errors.Add("expr", "invalid cron expression")
	}

	data := struct {
		Form *form.Form
	}{
		Form: f,
	}

	if !f.Valid() {
		app.render(w, r, files, data)
		return
	}

	expr := f.Get("expr")
	project := session.Account.Project
	schedule := core.NewSchedule(name, expr, project)
	err = app.storage.Schedule.Create(&schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Printf("account %s schedule %s create\n", session.Account.Email, schedule.ID)
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

	schedule, err := app.storage.Schedule.Read(id)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.storage.Schedule.Delete(schedule)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.infoLog.Printf("account %s schedule %s delete\n", session.Account.Email, schedule.ID)
	http.Redirect(w, r, "/schedule", http.StatusSeeOther)
}
