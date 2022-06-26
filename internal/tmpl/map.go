package tmpl

import (
	"html/template"
	"io/fs"
)

type Map map[string]*template.Template

// Based on:
// Let's Go - Chapter 5.3 (Alex Edwards)
func NewMap(dir fs.FS) (Map, error) {
	cache := make(map[string]*template.Template)

	// TODO: can this be done with one glob?
	pages, err := fs.Glob(dir, "*.page.html")
	if err != nil {
		return nil, err
	}

	subPages, err := fs.Glob(dir, "**/*.page.html")
	if err != nil {
		return nil, err
	}

	pages = append(pages, subPages...)

	// Create a unique template set for each page.
	for _, page := range pages {
		// Parse layouts into the template set.
		ts, err := template.ParseFS(dir, "*.layout.html")
		if err != nil {
			return nil, err
		}

		// Parse page into the template set.
		ts, err = ts.ParseFS(dir, page)
		if err != nil {
			return nil, err
		}

		cache[page] = ts
	}

	return cache, nil
}
