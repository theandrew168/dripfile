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

func NewApplication(
	logger *slog.Logger,
) *Application {
	app := Application{
		logger: logger,
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

	mux.HandleFunc("/itinerary", app.handleItineraryCreate(), "POST")
	mux.HandleFunc("/itinerary", app.handleItineraryList(), "GET")
	mux.HandleFunc("/itinerary/:id", app.handleItineraryRead(), "GET")
	mux.HandleFunc("/itinerary/:id", app.handleItineraryDelete(), "DELETE")
	mux.HandleFunc("/itinerary/:id/run", app.handleItineraryRun(), "POST")

	return mux
}
