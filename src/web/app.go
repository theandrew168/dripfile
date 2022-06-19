package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/src/config"
	"github.com/theandrew168/dripfile/src/jsonlog"
	"github.com/theandrew168/dripfile/src/secret"
	"github.com/theandrew168/dripfile/src/storage"
	"github.com/theandrew168/dripfile/src/stripe"
)

//go:embed static/img/logo-white.svg
var Favicon []byte

//go:embed static/etc/robots.txt
var Robots []byte

//go:embed static
var staticFS embed.FS

//go:embed template
var templateFS embed.FS

type Application struct {
	static   fs.FS
	template fs.FS

	logger  *jsonlog.Logger
	cfg     config.Config
	storage *storage.Storage
	queue   *asynq.Client
	box     *secret.Box
	billing stripe.Billing
}

func NewApplication(
	logger *jsonlog.Logger,
	cfg config.Config,
	storage *storage.Storage,
	queue *asynq.Client,
	box *secret.Box,
	billing stripe.Billing,
) *Application {
	var static fs.FS
	var template fs.FS
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		static = os.DirFS("./src/web/static/")
		template = os.DirFS("./src/web/template/")
	} else {
		// else use the embedded template FS
		var err error
		static, err = fs.Sub(staticFS, "static")
		if err != nil {
			panic(err)
		}
		template, err = fs.Sub(templateFS, "template")
		if err != nil {
			panic(err)
		}
	}

	app := Application{
		static:   static,
		template: template,

		logger:  logger,
		cfg:     cfg,
		storage: storage,
		queue:   queue,
		box:     box,
		billing: billing,
	}

	return &app
}

func (app *Application) Handler(api http.Handler) http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(app.recoverPanic)

	// healthcheck endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// top-level static files
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Write(Favicon)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write(Robots)
	})

	// static files
	staticServer := http.FileServer(http.FS(app.static))
	mux.Handle("/static/...", http.StripPrefix("/static", staticServer))

	// serve API routes under /api/v1
	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", api))
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusMovedPermanently)
	})

	// primary web app (last due to being a top-level catch-all)

	// landing page, visible to anyone
	mux.HandleFunc("/", app.handleIndex, "GET")

	// register / login / logout pages, visible to anyone
	mux.HandleFunc("/register", app.handleRegister, "GET")
	mux.Handle("/register", app.parseFormFunc(app.handleRegisterForm), "POST")
	mux.HandleFunc("/login", app.handleLogin, "GET")
	mux.Handle("/login", app.parseFormFunc(app.handleLoginForm), "POST")
	mux.Handle("/logout", app.parseFormFunc(app.handleLogoutForm), "POST")

	// app pages, visible only to authenticated users
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requireAuth)

		// only requires auth because billing might not be setup yet
		mux.HandleFunc("/billing/setup", app.handleBillingSetup, "GET")
		mux.HandleFunc("/billing/checkout", app.handleBillingCheckout, "GET")
		mux.HandleFunc("/billing/success", app.handleBillingSuccess, "GET")
		mux.HandleFunc("/billing/cancel", app.handleBillingCancel, "GET")

		// only requires auth so that new / old accounts can be managed w/o billing
		mux.HandleFunc("/account", app.handleAccountRead, "GET")
		mux.Handle("/account/delete", app.parseFormFunc(app.handleAccountDeleteForm), "POST")

		// these routes require auth AND billing
		mux.Group(func(mux *flow.Mux) {
			mux.Use(app.requireBillingSetup)

			mux.HandleFunc("/dashboard", app.handleDashboard, "GET")

			mux.HandleFunc("/transfer", app.handleTransferList, "GET")
			mux.HandleFunc("/transfer/create", app.handleTransferCreate, "GET")
			mux.Handle("/transfer/create", app.parseFormFunc(app.handleTransferCreateForm), "POST")
			mux.Handle("/transfer/delete", app.parseFormFunc(app.handleTransferDeleteForm), "POST")
			mux.Handle("/transfer/run", app.parseFormFunc(app.handleTransferRunForm), "POST")
			mux.HandleFunc("/transfer/:id", app.handleTransferRead, "GET")

			mux.HandleFunc("/location", app.handleLocationList, "GET")
			mux.HandleFunc("/location/create", app.handleLocationCreate, "GET")
			mux.Handle("/location/create", app.parseFormFunc(app.handleLocationCreateForm), "POST")
			mux.Handle("/location/delete", app.parseFormFunc(app.handleLocationDeleteForm), "POST")
			mux.HandleFunc("/location/:id", app.handleLocationRead, "GET")

			mux.HandleFunc("/schedule", app.handleScheduleList, "GET")
			mux.HandleFunc("/schedule/create", app.handleScheduleCreate, "GET")
			mux.Handle("/schedule/create", app.parseFormFunc(app.handleScheduleCreateForm), "POST")
			mux.Handle("/schedule/delete", app.parseFormFunc(app.handleScheduleDeleteForm), "POST")
			mux.HandleFunc("/schedule/:id", app.handleScheduleRead, "GET")

			mux.HandleFunc("/history", app.handleHistoryList, "GET")
		})
	})

	return mux
}
