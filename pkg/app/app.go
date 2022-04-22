package app

import (
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/secret"
	"github.com/theandrew168/dripfile/pkg/static"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/stripe"
	"github.com/theandrew168/dripfile/pkg/task"
	"github.com/theandrew168/dripfile/pkg/web"
)

// runtime deps (migrate, scheduler, worker, web):
// config
// logger
// storage
// queue
// box
// stripe
// mailer

// create the main application
func New(
	cfg config.Config,
	logger *jsonlog.Logger,
	storage *storage.Storage,
	queue *task.Queue,
	box *secret.Box,
	billing stripe.Billing,
) http.Handler {
	mux := flow.New()

	// handle top-level special cases
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	})
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(static.Favicon)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(static.Robots)
	})

	// static files app
	staticApp := static.NewApplication()
	mux.Handle("/static/...", http.StripPrefix("/static", staticApp.Router()))

	// rest api app
	//	apiApp := api.NewApplication(cfg, storage, logger)
	//	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", apiApp.Router()))

	// primary web app (last due to being a top-level catch-all)
	webApp := web.NewApplication(cfg, logger, storage, queue, box, billing)
	mux.Handle("/...", webApp.Router())

	return mux
}
