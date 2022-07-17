package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"
	"github.com/go-playground/form/v4"
	"github.com/hibiken/asynq"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/theandrew168/dripfile/internal/jsonlog"
	"github.com/theandrew168/dripfile/internal/secret"
	"github.com/theandrew168/dripfile/internal/storage"
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
	template *TemplateCache
	decoder  *form.Decoder

	logger *jsonlog.Logger
	store  *storage.Storage
	box    *secret.Box

	asynqClient *asynq.Client
}

func NewApplication(
	logger *jsonlog.Logger,
	store *storage.Storage,
	box *secret.Box,
	asynqClient *asynq.Client,
) *Application {
	static, err := fs.Sub(staticFS, "static")
	if err != nil {
		panic(err)
	}

	var template *TemplateCache
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		dir := os.DirFS("./internal/web/template/")
		template, err = NewTemplateCache(dir, true)
		if err != nil {
			panic(err)
		}
	} else {
		// else use the embedded template FS
		dir, err := fs.Sub(templateFS, "template")
		if err != nil {
			panic(err)
		}
		template, err = NewTemplateCache(dir, false)
		if err != nil {
			panic(err)
		}
	}

	// use a single instance of Decoder, it caches struct info
	decoder := form.NewDecoder()

	app := Application{
		static:   static,
		template: template,
		decoder:  decoder,

		logger: logger,
		store:  store,
		box:    box,

		asynqClient: asynqClient,
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

func (app *Application) Handler(api http.Handler) http.Handler {
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
