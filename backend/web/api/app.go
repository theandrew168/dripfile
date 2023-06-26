package api

import (
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger

	locationStorage location.Storage
	transferStorage transfer.Storage
	historyStorage  history.Storage

	transferService transfer.Service
}

func NewApplication(
	logger *slog.Logger,
	locationStorage location.Storage,
	transferStorage transfer.Storage,
	historyStorage history.Storage,
	transferService transfer.Service,
) *Application {
	app := Application{
		logger: logger,

		locationStorage: locationStorage,
		transferStorage: transferStorage,
		historyStorage:  historyStorage,

		transferService: transferService,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(middleware.RecoverPanic)
	mux.Use(middleware.SecureHeaders)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndex, "GET")

	mux.HandleFunc("/locations", app.handleListLocations, "GET")

	return mux
}
