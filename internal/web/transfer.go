package web

import (
	"net/http"
)

func (app *Application) handleTransferReadMany(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Category string
	}{
		Category: "transfer",
	}

	files := []string{
		"base.layout.html",
		"app.layout.html",
		"transfer/read_many.page.html",
	}

	app.render(w, r, files, data)
}
