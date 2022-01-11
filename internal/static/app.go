package static

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/klauspost/compress/gzhttp"

	"github.com/theandrew168/dripfile/internal/config"
)

//go:embed static/etc/favicon.webp
var Favicon []byte

//go:embed static/etc/robots.txt
var Robots []byte

//go:embed static
var staticFS embed.FS

type Application struct {
	cfg    config.Config
	static fs.FS
	logger *log.Logger
}

func NewApplication(cfg config.Config, logger *log.Logger) *Application {
	static, _ := fs.Sub(staticFS, "static")
	app := Application{
		cfg:    cfg,
		static: static,
		logger: logger,
	}

	return &app
}

func (app *Application) Router() http.Handler {
	// setup automatic compression handler for static files
	staticServer := http.FileServer(http.FS(app.static))
	gzipStaticServer := gzhttp.GzipHandler(staticServer)
	return gzipStaticServer
}
