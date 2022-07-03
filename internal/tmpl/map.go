package tmpl

import (
	"html/template"
	"io/fs"
	"strings"
)

type Map map[string]*template.Template

// Based on:
// Let's Go - Chapter 5.3 (Alex Edwards)
func NewMap(dir fs.FS) (Map, error) {
	cache := make(map[string]*template.Template)

	appPages, err := listTemplates(dir, "app")
	if err != nil {
		return nil, err
	}

	sitePages, err := listTemplates(dir, "site")
	if err != nil {
		return nil, err
	}

	errorPages, err := listTemplates(dir, "error")
	if err != nil {
		return nil, err
	}

	var pages []string
	pages = append(pages, appPages...)
	pages = append(pages, sitePages...)
	pages = append(pages, errorPages...)

	// Create a unique template set for each page.
	for _, page := range pages {
		ts, err := parseTemplate(dir, page)
		if err != nil {
			return nil, err
		}

		cache[page] = ts
	}

	return cache, nil
}

func listTemplates(dir fs.FS, root string) ([]string, error) {
	var templates []string
	err := fs.WalkDir(dir, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".html") {
			templates = append(templates, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return templates, nil
}

func parseTemplate(dir fs.FS, page string) (*template.Template, error) {
	// base layout
	ts, err := template.ParseFS(dir, "layout/base.html")
	if err != nil {
		return nil, err
	}

	// sub-layouts (if necessary)
	if strings.HasPrefix(page, "app/") {
		ts, err = ts.ParseFS(dir, "layout/app.html")
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(page, "site/") {
		ts, err = ts.ParseFS(dir, "layout/site.html")
		if err != nil {
			return nil, err
		}
	}

	// partials
	partials, err := listTemplates(dir, "partial")
	if err != nil {
		return nil, err
	}

	ts, err = ts.ParseFS(dir, partials...)
	if err != nil {
		return nil, err
	}

	// page
	ts, err = ts.ParseFS(dir, page)
	if err != nil {
		return nil, err
	}

	return ts, nil
}
