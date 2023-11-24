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

	mux.HandleFunc("/locations", app.handleLocationCreate(), "POST")
	mux.HandleFunc("/locations", app.handleLocationList(), "GET")
	mux.HandleFunc("/locations/:id", app.handleLocationRead(), "GET")
	mux.HandleFunc("/locations/:id", app.handleLocationDelete(), "DELETE")

	mux.HandleFunc("/itineraries", app.handleItineraryCreate(), "POST")
	mux.HandleFunc("/itineraries", app.handleItineraryList(), "GET")
	mux.HandleFunc("/itineraries/:id", app.handleItineraryRead(), "GET")
	mux.HandleFunc("/itineraries/:id", app.handleItineraryDelete(), "DELETE")

	mux.HandleFunc("/transfers", app.handleTransferCreate(), "POST")
	mux.HandleFunc("/transfers", app.handleTransferList(), "GET")
	mux.HandleFunc("/transfers/:id", app.handleTransferRead(), "GET")

	return mux
}
