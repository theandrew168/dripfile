package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/core"
)

//go:embed templates
var templatesFS embed.FS

type Application struct {
	templates fs.FS

	cfg     config.Config
	storage core.Storage
	logger  *log.Logger
}

func NewApplication(cfg config.Config, storage core.Storage, logger *log.Logger) *Application {
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

		cfg:     cfg,
		storage: storage,
		logger:  logger,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)
	mux.HandleFunc("/", app.handleIndex, "GET")
	return mux
}
