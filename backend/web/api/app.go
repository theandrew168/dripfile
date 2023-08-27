package api

import (
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/history"
	"github.com/theandrew168/dripfile/backend/location"
	"github.com/theandrew168/dripfile/backend/transfer"
	transferService "github.com/theandrew168/dripfile/backend/transfer/service"
	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger

	locationStorage location.Repository
	transferStorage transfer.Repository
	historyStorage  history.Repository

	transferService transferService.Service
}

func NewApplication(
	logger *slog.Logger,
	locationStorage location.Repository,
	transferStorage transfer.Repository,
	historyStorage history.Repository,
	transferService transferService.Service,
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

	mux.HandleFunc("/locations", app.handleLocationList, "GET")
	mux.HandleFunc("/locations/:id", app.handleLocationRead, "GET")

	return mux
}
