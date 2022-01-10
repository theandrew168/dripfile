package web

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/alexedwards/flow"
	"github.com/klauspost/compress/gzhttp"

	"github.com/theandrew168/dripfile/internal/config"
	"github.com/theandrew168/dripfile/internal/core"
)

//go:embed static/img/logo.webp
var logo []byte

//go:embed static
var staticFS embed.FS

//go:embed templates
var templatesFS embed.FS

type Application struct {
	cfg config.Config

	static    fs.FS
	templates fs.FS
	storage   core.Storage
	logger    *log.Logger
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

	static, _ := fs.Sub(staticFS, "static")
	app := Application{
		cfg: cfg,

		static:    static,
		templates: templates,
		storage:   storage,
		logger:    logger,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	// setup http.Handler for static files
	static, _ := fs.Sub(staticFS, "static")
	staticServer := http.FileServer(http.FS(static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)

	mux := flow.New()
	mux.HandleFunc("/", app.handleIndex, "GET")

	// static files (compressed) and favicon (special case)
	mux.Handle("/static/...", http.StripPrefix("/static", gzipStaticServer))
	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/webp")
		w.Write(logo)
	})

	return mux
}
