package web

import (
	"io/fs"
	"net/http"
	"os"

	"github.com/alexedwards/flow"
	"github.com/klauspost/compress/gzhttp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
	transferService "github.com/theandrew168/dripfile/backend/transfer/service"
	"github.com/theandrew168/dripfile/backend/web/api"
	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger
	public fs.FS

	locationStorage location.Repository
	transferStorage transfer.Repository
	historyStorage  history.Repository

	transferService transferService.Service
}

func NewApplication(
	logger *slog.Logger,
	publicFS fs.FS,
	locationStorage location.Repository,
	transferStorage transfer.Repository,
	historyStorage history.Repository,
	transferService transferService.Service,
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

		locationStorage: locationStorage,
		transferStorage: transferStorage,
		historyStorage:  historyStorage,

		transferService: transferService,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	mux.Use(middleware.RecoverPanic)

	// healthcheck endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// REST API routes
	apiV1 := api.NewApplication(
		app.logger,
		app.locationStorage,
		app.transferStorage,
		app.historyStorage,
		app.transferService,
	)
	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", apiV1.Handler()))

	// public files to be served (and auto-compressed)
	public := gzhttp.GzipHandler(http.FileServer(http.FS(app.public)))
	mux.Handle("/", public)
	mux.Handle("/index.html", public)
	mux.Handle("/index.js", public)
	mux.Handle("/index.css", public)
	mux.Handle("/robots.txt", public)
	mux.Handle("/favicon.ico", public)
	mux.Handle("/static/...", public)

	// all other routes should return the index page
	// so that the frontend router can take over
	mux.Handle("/...", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		index, err := fs.ReadFile(app.public, "index.html")
		if err != nil {
			panic(err)
		}

		w.Write(index)
	}))

	return mux
}
