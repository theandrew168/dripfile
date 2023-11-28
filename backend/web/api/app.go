package api

import (
	"log/slog"
	"net/http"

	"github.com/alexedwards/flow"

	"github.com/theandrew168/dripfile/backend/repository"
	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger
	repo   *repository.Repository
}

func NewApplication(
	logger *slog.Logger,
	repo *repository.Repository,
) *Application {
	app := Application{
		logger: logger,
		repo:   repo,
	}
	return &app
}

func (app *Application) Handler() http.Handler {
	mux := flow.New()
	mux.NotFound = http.HandlerFunc(app.notFoundResponse)
	mux.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	mux.Use(middleware.RecoverPanic)
	mux.Use(middleware.EnableCORS)

	mux.HandleFunc("/", app.handleIndex(), "GET")

	mux.HandleFunc("/location", app.handleLocationCreate(), "POST")
	mux.HandleFunc("/location", app.handleLocationList(), "GET")
	mux.HandleFunc("/location/:id", app.handleLocationRead(), "GET")
	mux.HandleFunc("/location/:id", app.handleLocationDelete(), "DELETE")
	mux.HandleFunc("/location/:id/ping", app.handleLocationPing(), "POST")

	mux.HandleFunc("/itinerary", app.handleItineraryCreate(), "POST")
	mux.HandleFunc("/itinerary", app.handleItineraryList(), "GET")
	mux.HandleFunc("/itinerary/:id", app.handleItineraryRead(), "GET")
	mux.HandleFunc("/itinerary/:id", app.handleItineraryDelete(), "DELETE")

	mux.HandleFunc("/transfer", app.handleTransferCreate(), "POST")
	mux.HandleFunc("/transfer", app.handleTransferList(), "GET")
	mux.HandleFunc("/transfer/:id", app.handleTransferRead(), "GET")

	return mux
}
