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

//go:embed templates
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
		templates = os.DirFS("./internal/web/templates/")
	} else {
		// else use the embedded templates dir
		templates, _ = fs.Sub(templatesFS, "templates")
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

	// landing page, visible to anyone
	mux.HandleFunc("/", app.handleIndex, "GET")

	// login / register pages, visible to anyone
	mux.HandleFunc("/login", app.handleLogin, "GET")
	mux.HandleFunc("/login", app.handleLoginForm, "POST")
	mux.HandleFunc("/register", app.handleRegister, "GET")
	mux.HandleFunc("/register", app.handleRegisterForm, "POST")

	// app pages, visible only to authenticated users
	mux.HandleFunc("/dashboard", app.handleDashboard, "GET")

	return mux
}
