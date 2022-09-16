package web

import (
	"net/http"

	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/service"
	"github.com/theandrew168/dripfile/internal/storage"
	"github.com/theandrew168/dripfile/internal/task"
	"github.com/theandrew168/dripfile/internal/view"
)

type Application struct {
	decoder *form.Decoder

	static http.Handler
	api    http.Handler
	logger *jsonlog.Logger
	view   *view.View
	srvc   *service.Service
	store  *storage.Storage
	queue  *task.Queue
	box    *secret.Box
}

func NewApplication(
	api http.Handler,
	static http.Handler,
	logger *jsonlog.Logger,
	view *view.View,
	srvc *service.Service,
	store *storage.Storage,
	queue *task.Queue,
	box *secret.Box,
) *Application {
	// use a single instance of Decoder (it caches struct info)
	decoder := form.NewDecoder()

	app := Application{
		decoder: decoder,

		api:    api,
		static: static,
		logger: logger,
		view:   view,
		srvc:   srvc,
		store:  store,
		queue:  queue,
		box:    box,
	}

	return &app
}

// Redirects:
// 303 See Other         - for GETs after POSTs (like a login / register form)
// 302 Found             - all other temporary redirects
// 301 Moved Permanently - permanent redirects

// Route Handler Naming Ideas:
//
// basic page handlers:
// GET - handleIndex
// GET - handleDashboard
//
// basic page w/ form handlers:
// GET  - handleLogin
// POST - handleLoginForm
//
// CRUD handlers:
// C POST   - handleCreateFoo[Form]
// R GET    - handleReadFoo[s]
// U PUT    - handleUpdateFoo[Form]
// D DELETE - handleDeleteFoo[Form]

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(app.recoverPanic)
	mux.Use(app.setSecureHeaders)
	mux.Use(app.limitRequestSize)

	// healthcheck endpoint
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("pong\n"))
	}, "GET")

	// prometheus metrics
	mux.Handle("/metrics", promhttp.Handler(), "GET")

	// static files (and top-level redirects)
	mux.Handle("/static/...", http.StripPrefix("/static", app.static))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/img/logo-white.svg", http.StatusMovedPermanently)
	})
	mux.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/static/etc/robots.txt", http.StatusMovedPermanently)
	})

	// serve API routes under /api/v1
	mux.Handle("/api/v1/...", http.StripPrefix("/api/v1", app.api))
	mux.HandleFunc("/api/v1", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/v1/", http.StatusMovedPermanently)
	})

	// primary web app (last due to being a top-level catch-all)

	// landing page, visible to anyone
	mux.HandleFunc("/", app.handleIndex, "GET")

	// register / login / logout pages, visible to anyone
	mux.HandleFunc("/register", app.handleRegister, "GET")
	mux.HandleFunc("/register", app.handleRegisterForm, "POST")
	mux.HandleFunc("/login", app.handleLogin, "GET")
	mux.HandleFunc("/login", app.handleLoginForm, "POST")
	mux.HandleFunc("/logout", app.handleLogoutForm, "POST")

	// app pages, visible only to authenticated users
	mux.Group(func(mux *flow.Mux) {
		mux.Use(app.requireAuth)

		mux.HandleFunc("/dashboard", app.handleDashboard, "GET")

		mux.HandleFunc("/account", app.handleAccountRead, "GET")
		mux.HandleFunc("/account/delete", app.handleAccountDeleteForm, "POST")

		mux.HandleFunc("/transfer", app.handleTransferList, "GET")
		mux.HandleFunc("/transfer/create", app.handleTransferCreate, "GET")
		mux.HandleFunc("/transfer/create", app.handleTransferCreateForm, "POST")
		mux.HandleFunc("/transfer/delete", app.handleTransferDeleteForm, "POST")
		mux.HandleFunc("/transfer/run", app.handleTransferRunForm, "POST")
		mux.HandleFunc("/transfer/:id", app.handleTransferRead, "GET")

		mux.HandleFunc("/location", app.handleLocationList, "GET")
		mux.HandleFunc("/location/create", app.handleLocationCreate, "GET")
		mux.HandleFunc("/location/create", app.handleLocationCreateForm, "POST")
		mux.HandleFunc("/location/delete", app.handleLocationDeleteForm, "POST")
		mux.HandleFunc("/location/:id", app.handleLocationRead, "GET")

		mux.HandleFunc("/schedule", app.handleScheduleList, "GET")
		mux.HandleFunc("/schedule/create", app.handleScheduleCreate, "GET")
		mux.HandleFunc("/schedule/create", app.handleScheduleCreateForm, "POST")
		mux.HandleFunc("/schedule/delete", app.handleScheduleDeleteForm, "POST")
		mux.HandleFunc("/schedule/:id", app.handleScheduleRead, "GET")

		mux.HandleFunc("/history", app.handleHistoryList, "GET")
	})

	return mux
}
