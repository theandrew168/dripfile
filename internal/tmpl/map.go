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

	appPages, err := listPages(dir, "app")
	if err != nil {
		return nil, err
	}

	sitePages, err := listPages(dir, "site")
	if err != nil {
		return nil, err
	}

	errorPages, err := listPages(dir, "error")
	if err != nil {
		return nil, err
	}

	var pages []string
	pages = append(pages, appPages...)
	pages = append(pages, sitePages...)
	pages = append(pages, errorPages...)

	// Create a unique template set for each page.
	for _, page := range pages {
		ts, err := readPage(dir, page)
		if err != nil {
			return nil, err
		}

		cache[page] = ts
	}

	return cache, nil
}

func listPages(dir fs.FS, root string) ([]string, error) {
	var pages []string
	err := fs.WalkDir(dir, root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if strings.HasSuffix(path, ".page.html") {
			pages = append(pages, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return pages, nil
}

func readPage(dir fs.FS, page string) (*template.Template, error) {
	// Parse base layout into the template set.
	ts, err := template.ParseFS(dir, "base.layout.html")
	if err != nil {
		return nil, err
	}

	// Parse sub-layouts if necessary
	if strings.HasPrefix(page, "app/") {
		ts, err = ts.ParseFS(dir, "app.layout.html")
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(page, "site/") {
		ts, err = ts.ParseFS(dir, "site.layout.html")
		if err != nil {
			return nil, err
		}
	}

	// Parse page into the template set.
	ts, err = ts.ParseFS(dir, page)
	if err != nil {
		return nil, err
	}

	return ts, nil
}
