package web

import (
	"net/http"
)

func (app *Application) handleBillingCheckout(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create checkout session
	billingID := session.Account.Project.BillingID
	sessionURL, err := app.paygate.CreateCheckoutSession(billingID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// redirect user to Stripe to provide payment info
	http.Redirect(w, r, sessionURL, http.StatusFound)
}

func (app *Application) handleBillingSuccess(w http.ResponseWriter, r *http.Request) {
	// TODO: get checkout session
	// TODO: get setup intent (expand into one call?)
	// TODO: store something? payment ID? just mark success?

	app.logger.Info("%s\n", r.URL)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleBillingCancel(w http.ResponseWriter, r *http.Request) {
	// TODO: redir to dashboard, middleware will catch missing
	// 	payment info and retry the checkout session?

	app.logger.Info("%s\n", r.URL)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
