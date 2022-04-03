package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/pkg/config"
	"github.com/theandrew168/dripfile/pkg/secret"
	"github.com/theandrew168/dripfile/pkg/storage"
	"github.com/theandrew168/dripfile/pkg/stripe"
	"github.com/theandrew168/dripfile/pkg/task"
)

//go:embed template
var templateFS embed.FS

type Application struct {
	templates fs.FS

	cfg      config.Config
	box      *secret.Box
	storage  *storage.Storage
	queue    *task.Queue
	stripe   stripe.Interface
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewApplication(
	cfg config.Config,
	box *secret.Box,
	storage *storage.Storage,
	queue *task.Queue,
	stripe stripe.Interface,
	infoLog *log.Logger,
	errorLog *log.Logger,
) *Application {
	var templates fs.FS
	if strings.HasPrefix(os.Getenv("ENV"), "dev") {
		// reload templates from filesystem if var ENV starts with "dev"
		// NOTE: os.DirFS is rooted from where the app is ran, not this file
		templates = os.DirFS("./pkg/web/template/")
	} else {
		// else use the embedded template FS
		var err error
		templates, err = fs.Sub(templateFS, "template")
		if err != nil {
			panic(err)
		}
	}

	app := Application{
		templates: templates,

		cfg:      cfg,
		box:      box,
		storage:  storage,
		queue:    queue,
		stripe:   stripe,
		infoLog:  infoLog,
		errorLog: errorLog,
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

		// TODO: move into separate stripe/app.go?
		mux.HandleFunc("/stripe/payment", app.handleStripePayment, "GET")
		mux.HandleFunc("/stripe/checkout", app.handleStripeCheckout, "GET")
		mux.HandleFunc("/stripe/success", app.handleStripeSuccess, "GET")
		mux.HandleFunc("/stripe/cancel", app.handleStripeCancel, "GET")

		mux.HandleFunc("/account", app.handleAccountRead, "GET")
		mux.Handle("/account/delete", app.parseFormFunc(app.handleAccountDeleteForm), "POST")

		mux.Group(func(mux *flow.Mux) {
			mux.Use(app.requirePaymentInfo)

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

			mux.HandleFunc("/history", app.handleHistoryList, "GET")
		})
	})

	return mux
}
