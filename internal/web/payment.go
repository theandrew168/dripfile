package web

import (
	"net/http"
)

func (app *Application) handleCheckout(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"payment/checkout.page.html",
	}

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create setup intent
	billingID := session.Account.Project.BillingID
	clientSecret, err := app.billing.CreateSetupIntent(billingID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := struct {
		ClientSecret    string
		StripePublicKey string
		ReturnURL       string
	}{
		ClientSecret:    clientSecret,
		StripePublicKey: app.cfg.StripePublicKey,
		ReturnURL:       "http://localhost:5000/status",
	}

	app.render(w, r, files, data)
}

func (app *Application) handleStatus(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"payment/status.page.html",
	}

	data := struct {
		StripePublicKey string
	}{
		StripePublicKey: app.cfg.StripePublicKey,
	}

	app.render(w, r, files, data)
}
