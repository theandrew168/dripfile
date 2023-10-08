package api

import (
	"net/http"

	"github.com/alexedwards/flow"
	"golang.org/x/exp/slog"

	"github.com/theandrew168/dripfile/backend/service"
	"github.com/theandrew168/dripfile/backend/web/middleware"
)

type Application struct {
	logger *slog.Logger
	srvc   *service.Service
}

func NewApplication(
	logger *slog.Logger,
	srvc *service.Service,
) *Application {
	app := Application{
		logger: logger,
		srvc:   srvc,
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

	mux.HandleFunc("/transfer", app.handleTransferCreate(), "POST")
	mux.HandleFunc("/transfer", app.handleTransferList(), "GET")
	mux.HandleFunc("/transfer/:id", app.handleTransferRead(), "GET")

	return mux
}
