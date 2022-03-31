// https://stripe.com/docs/payments/save-and-reuse?platform=checkout
package billing

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

type stripeGateway struct {
	client     *client.API
	publicKey  string
	secretKey  string
	successURL string
	cancelURL  string
}

func NewStripeGateway(publicKey, secretKey, successURL, cancelURL string) PaymentGateway {
	sc := &client.API{}
	sc.Init(secretKey, nil)

	g := stripeGateway{
		client:     sc,
		publicKey:  publicKey,
		secretKey:  secretKey,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
	return &g
}

func (g *stripeGateway) CreateCustomer(email string) (string, error) {
	params := stripe.CustomerParams{
		Email: stripe.String(email),
	}

	customer, err := g.client.Customers.New(&params)
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func (g *stripeGateway) CreateCheckoutSession(billingID string) (string, error) {
	params := stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSetup)),
		Customer:   stripe.String(billingID),
		SuccessURL: stripe.String(g.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(g.cancelURL),
	}

	session, err := g.client.CheckoutSessions.New(&params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}
