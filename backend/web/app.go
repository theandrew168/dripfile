package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/alexedwards/flow"
	"github.com/klauspost/compress/gzhttp"
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
	var public fs.FS
	if os.Getenv("DEBUG") != "" {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		public = os.DirFS("./public/")
	} else {
		// else use the embedded template dir
		public, _ = fs.Sub(publicFS, "public")
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
	mux.Handle("/...", gzhttp.GzipHandler(public))

	return mux
}
