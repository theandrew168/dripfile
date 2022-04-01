// https://stripe.com/docs/payments/save-and-reuse?platform=checkout
package billing

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

const stripePriceID = "price_1KhKCdFGWxTaRTVh9iXdyx0M"

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

func (g *stripeGateway) CreateCheckoutSession(customerID string) (string, error) {
	params := stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSetup)),
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String(g.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(g.cancelURL),
	}

	session, err := g.client.CheckoutSessions.New(&params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func (g *stripeGateway) CreateSubscription(customerID string) (string, error) {
	params := stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(stripePriceID),
			},
		},
	}

	subscription, err := g.client.Subscriptions.New(&params)
	if err != nil {
		return "", err
	}

	return subscription.Items.Data[0].ID, nil
}
