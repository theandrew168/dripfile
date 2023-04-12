package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/validator"
)

type ScheduleDeleteForm struct {
	validator.Validator `form:"-"`

	ScheduleID string `form:"ScheduleID"`
}

type ScheduleCreateForm struct {
	validator.Validator `form:"-"`

	Expr string `form:"Expr"`
}

type ScheduleCreateParams struct {
	Form ScheduleCreateForm
}

func (v *View) ScheduleCreate(w io.Writer, p ScheduleCreateParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"schedule/create.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}

type ScheduleReadParams struct {
	Schedule model.Schedule
}

func (v *View) ScheduleRead(w io.Writer, p ScheduleReadParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"schedule/read.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}

type ScheduleListParams struct {
	Schedules []model.Schedule
}

func (v *View) ScheduleList(w io.Writer, p ScheduleListParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"schedule/list.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}
