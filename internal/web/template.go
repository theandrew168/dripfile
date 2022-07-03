package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"strings"
)

type TemplateCache struct {
	dir    fs.FS
	cache  map[string]*template.Template
	reload bool
}

// Based on:
// Let's Go - Chapter 5.3 (Alex Edwards)
func NewTemplateCache(dir fs.FS, reload bool) (*TemplateCache, error) {
	tc := TemplateCache{
		dir:    dir,
		cache:  make(map[string]*template.Template),
		reload: reload,
	}

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
		t, err := parseTemplate(dir, page)
		if err != nil {
			return nil, err
		}

		tc.cache[page] = t
	}

	return &tc, nil
}

func (tc *TemplateCache) Get(page string) (*template.Template, error) {
	if tc.reload {
		t, err := parseTemplate(tc.dir, page)
		if err != nil {
			return nil, err
		}

		return t, nil
	}

	t, ok := tc.cache[page]
	if !ok {
		err := fmt.Errorf("web: template does not exist: %s", page)
		return nil, err
	}

	return t, nil
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
	t, err := template.ParseFS(dir, "layout/base.html")
	if err != nil {
		return nil, err
	}

	// sub-layouts (if necessary)
	if strings.HasPrefix(page, "app/") {
		t, err = t.ParseFS(dir, "layout/app.html")
		if err != nil {
			return nil, err
		}
	} else if strings.HasPrefix(page, "site/") {
		t, err = t.ParseFS(dir, "layout/site.html")
		if err != nil {
			return nil, err
		}
	}

	// partials
	partials, err := listTemplates(dir, "partial")
	if err != nil {
		return nil, err
	}

	t, err = t.ParseFS(dir, partials...)
	if err != nil {
		return nil, err
	}

	// page
	t, err = t.ParseFS(dir, page)
	if err != nil {
		return nil, err
	}

	return t, nil
}
