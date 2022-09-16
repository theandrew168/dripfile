package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/validator"
)

type TransferRunForm struct {
	validator.Validator `form:"-"`

	TransferID string `form:"TransferID"`
}

type TransferDeleteForm struct {
	validator.Validator `form:"-"`

	TransferID string `form:"TransferID"`
}

type TransferCreateForm struct {
	validator.Validator `form:"-"`

	Pattern    string `form:"Pattern"`
	SrcID      string `form:"SrcID"`
	DstID      string `form:"DstID"`
	ScheduleID string `form:"ScheduleID"`
}

type TransferCreateParams struct {
	Form TransferCreateForm

	Locations []model.Location `form:"Locations"`
	Schedules []model.Schedule `form:"Schedules"`
}

func (v *View) TransferCreate(w io.Writer, p TransferCreateParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"transfer/create.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}

type TransferReadParams struct {
	Transfer model.Transfer
}

func (v *View) TransferRead(w io.Writer, p TransferReadParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"transfer/read.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}

type TransferListParams struct {
	Transfers []model.Transfer
}

func (v *View) TransferList(w io.Writer, p TransferListParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"transfer/list.html",
	}
	tmpl := v.r.Read(patterns...)
	return tmpl.Execute(w, p)
}
