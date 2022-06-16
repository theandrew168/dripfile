package web

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/backend/jsonlog"
)

//go:embed public
var publicFS embed.FS

type Application struct {
	public fs.FS

	logger *jsonlog.Logger
}

func NewApplication(logger *jsonlog.Logger) *Application {
	public, err := fs.Sub(publicFS, "public")
	if err != nil {
		panic(err)
	}

	app := Application{
		public: public,

		logger: logger,
	}
	return &app
}

func (app *Application) Handler(api http.Handler) http.Handler {
	mux := flow.New()
	// TODO: web 404 / 405 pages

	// healthcheck endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// serve API routes under /api/v1
	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", api))
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusMovedPermanently)
	})

	// default to serving public files (svelte app, favicon, robots.txt, etc)
	mux.Handle("/...", http.FileServer(http.FS(app.public)))

	return mux
}
