package template

import (
	"html/template"
	"io/fs"
	"path"
	"strings"
)

type Reader struct {
	cache map[string]*template.Template

	files  fs.FS
	reload bool
}

func NewReader(files fs.FS, reload bool) *Reader {
	r := Reader{
		cache: make(map[string]*template.Template),

		files:  files,
		reload: reload,
	}
	return &r
}

func (r *Reader) Read(patterns ...string) *template.Template {
	if len(patterns) == 0 {
		return nil
	}

	// join the patterns to determine cache key
	key := strings.Join(patterns, ",")

	// load from cache if not doing dynamic reloads
	if !r.reload {
		tmpl, ok := r.cache[key]
		if ok {
			return tmpl
		}
	}

	// first pattern is the name / layout (needs to be a base path)
	// https://pkg.go.dev/text/template#Template.ParseFiles
	name := path.Base(patterns[0])

	tmpl, err := template.New(name).ParseFS(r.files, patterns...)
	if err != nil {
		panic(err)
	}

	r.cache[key] = tmpl
	return tmpl
}
