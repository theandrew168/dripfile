package template

import (
	"html/template"
	"io/fs"
	"path"
	"strings"
)

type Template struct {
	cache map[string]*template.Template

	files fs.FS
	reload bool
}

func New(files fs.FS, reload bool) Template {
	cache := make(map[string]*template.Template)

	t := Template{
		cache: cache,

		files: files,
		reload: reload,
	}
	return t
}

func (t *Template) Parse(patterns ...string) *template.Template {
	l := len(patterns)
	if l == 0 {
		return nil
	}

	// join the patterns to determine cache key
	key := strings.Join(patterns, ",")

	// load from cache if not doing dynamic reloads
	if !t.reload {
		tmpl, ok := t.cache[key]
		if ok {
			return tmpl
		}
	}

	// first pattern is the name / layout (needs to be a base path)
	// https://pkg.go.dev/text/template#Template.ParseFiles
	name := path.Base(patterns[0])

	tmpl, err := template.New(name).ParseFS(t.files, patterns...)
	if err != nil {
		panic(err)
	}

	t.cache[key] = tmpl
	return tmpl
}
