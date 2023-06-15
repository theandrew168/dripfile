package web

import (
	"io/fs"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"
)

type Application struct {
	logger *slog.Logger
	public fs.FS
}

func NewApplication(
	logger *slog.Logger,
	publicFS fs.FS,
) *Application {
	// drill-down one level so that the contents in public/
	// are served from the root URL of the app
	public, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic(err)
	}

	app := Application{
		logger: logger,
		public: public,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	// TODO: api 404 / 405 responses

	// healthcheck endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// TODO: /api/v1/ routes

	public := http.FileServer(http.FS(app.public))
	mux.Handle("/...", public)

	return mux
}
