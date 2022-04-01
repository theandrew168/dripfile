// https://stripe.com/docs/payments/save-and-reuse?platform=checkout
package stripe

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

const stripePriceID = "price_1KhKCdFGWxTaRTVh9iXdyx0M"

type stripeImpl struct {
	client     *client.API
	secretKey  string
	successURL string
	cancelURL  string
}

func New(secretKey, successURL, cancelURL string) Interface {
	sc := &client.API{}
	sc.Init(secretKey, nil)

	i := stripeImpl{
		client:     sc,
		secretKey:  secretKey,
		successURL: successURL,
		cancelURL:  cancelURL,
	}
	return &i
}

func (i *stripeImpl) CreateCustomer(email string) (string, error) {
	params := stripe.CustomerParams{
		Email: stripe.String(email),
	}

	customer, err := i.client.Customers.New(&params)
	if err != nil {
		return "", err
	}

	return customer.ID, nil
}

func (i *stripeImpl) CreateCheckoutSession(customerID string) (string, error) {
	params := stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSetup)),
		Customer:   stripe.String(customerID),
		SuccessURL: stripe.String(i.successURL + "?session_id={CHECKOUT_SESSION_ID}"),
		CancelURL:  stripe.String(i.cancelURL),
	}

	session, err := i.client.CheckoutSessions.New(&params)
	if err != nil {
		return "", err
	}

	return session.URL, nil
}

func (i *stripeImpl) CreateSubscription(customerID string) (string, error) {
	params := stripe.SubscriptionParams{
		Customer: stripe.String(customerID),
		Items: []*stripe.SubscriptionItemsParams{
			{
				Price: stripe.String(stripePriceID),
			},
		},
	}

	subscription, err := i.client.Subscriptions.New(&params)
	if err != nil {
		return "", err
	}

	return subscription.Items.Data[0].ID, nil
}
