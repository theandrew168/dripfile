package api

import (
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger
}

func NewApplication(logger *slog.Logger) *Application {
	app := Application{
		logger: logger,
	}
	return &app
}

func (app *Application) Router() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(middleware.RecoverPanic)
	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.HandleIndex, "GET")

	return mux
}
