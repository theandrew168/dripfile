package web

import (
	"fmt"
	"net/http"
)

// helper for rendering templates and handling potential errors
func (app *Application) render(w http.ResponseWriter, r *http.Request, page string, data interface{}) {
	ts, ok := app.tm[page]
	if !ok {
		err := fmt.Errorf("web: template does not exist: %s", page)
		app.serverErrorResponse(w, r, err)
		return
	}

	// always execute the base template
	w.Header().Set("Content-Type", "text/html")
	err := ts.ExecuteTemplate(w, "base", data)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
