package web

import (
	"embed"
	"io/fs"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/core"
	"github.com/theandrew168/dripfile/internal/log"
)

//go:embed template
var templatesFS embed.FS

type Application struct {
	templates fs.FS

	storage core.Storage
	logger  log.Logger
}

func NewApplication(storage core.Storage, logger log.Logger) *Application {
	var templates fs.FS
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./internal/web/template/")
	} else {
		// else use the embedded templates FS
		var err error
		templates, err = fs.Sub(templatesFS, "template")
		if err != nil {
			panic(err)
		}
	}

	app := Application{
		templates: templates,

		storage: storage,
		logger:  logger,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(app.recoverPanic)

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

		mux.HandleFunc("/dashboard", app.handleDashboard, "GET")

		mux.HandleFunc("/account", app.handleAccountRead, "GET")
		mux.Handle("/account/delete", app.parseFormFunc(app.handleAccountDeleteForm), "POST")

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

		mux.HandleFunc("/history", app.handleHistoryList, "GET")
	})

	return mux
}
