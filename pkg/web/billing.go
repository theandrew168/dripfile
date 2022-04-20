package web

import (
	"net/http"
)

func (app *Application) handleBillingSetup(w http.ResponseWriter, r *http.Request) {
	files := []string{
		"base.layout.html",
		"app.layout.html",
		"billing/setup.page.html",
	}

	app.render(w, r, files, nil)
}

func (app *Application) handleBillingCheckout(w http.ResponseWriter, r *http.Request) {
	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// create checkout session
	customerID := session.Account.Project.CustomerID
	sessionURL, err := app.stripe.CreateCheckoutSession(customerID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// redirect user to Billing to provide payment info
	http.Redirect(w, r, sessionURL, http.StatusFound)
}

func (app *Application) handleBillingSuccess(w http.ResponseWriter, r *http.Request) {
	// TODO: get checkout session
	// TODO: get setup intent (expand into one call?)
	// TODO: store something? payment ID? just mark success?
	// TODO: create subscription, store sub item ID in DB

	session, err := app.requestSession(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	customerID := session.Account.Project.CustomerID
	subscriptionItemID, err := app.stripe.CreateSubscription(customerID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	project := session.Account.Project
	project.SubscriptionItemID = subscriptionItemID
	err = app.storage.Project.Update(project)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (app *Application) handleBillingCancel(w http.ResponseWriter, r *http.Request) {
	// TODO: redir to dashboard, middleware will catch missing
	// 	payment info and retry the checkout session?

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}
