package web

import (
	"io"

	"github.com/theandrew168/dripfile/internal/model"
	"github.com/theandrew168/dripfile/internal/validator"
)

type LocationDeleteForm struct {
	validator.Validator `form:"-"`

	LocationID string `form:"LocationID"`
}

type LocationCreateForm struct {
	validator.Validator `form:"-"`

	Endpoint        string `form:"Endpoint"`
	BucketName      string `form:"BucketName"`
	AccessKeyID     string `form:"AccessKeyID"`
	SecretAccessKey string `form:"SecretAccessKey"`
}

type LocationCreateParams struct {
	Form LocationCreateForm
}

func (t *Template) LocationCreate(w io.Writer, p LocationCreateParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"location/create.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}

type LocationReadParams struct {
	Location model.Location
}

func (t *Template) LocationRead(w io.Writer, p LocationReadParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"location/read.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}

type LocationListParams struct {
	Locations []model.Location
}

func (t *Template) LocationList(w io.Writer, p LocationListParams) error {
	patterns := []string{
		"layout/base.html",
		"layout/app.html",
		"partial/*.html",
		"location/list.html",
	}
	tmpl := t.Parse(patterns...)
	return tmpl.Execute(w, p)
}
