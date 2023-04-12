package web

import (
	"io"
)

type Error400Params struct{}

func (html *HTML) Error400(w io.Writer, p Error400Params) error {
	patterns := []string{
		"layout/base.html",
		"error/400.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}

type Error404Params struct{}

func (html *HTML) Error404(w io.Writer, p Error404Params) error {
	patterns := []string{
		"layout/base.html",
		"error/404.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}

type Error405Params struct{}

func (html *HTML) Error405(w io.Writer, p Error405Params) error {
	patterns := []string{
		"layout/base.html",
		"error/405.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}

type Error500Params struct{}

func (html *HTML) Error500(w io.Writer, p Error500Params) error {
	patterns := []string{
		"layout/base.html",
		"error/500.html",
	}
	tmpl := html.reader.Read(patterns...)
	return tmpl.Execute(w, p)
}
