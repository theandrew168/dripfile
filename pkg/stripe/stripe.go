// https://stripe.com/docs/payments/save-and-reuse?platform=checkout
package stripe

import (
	"fmt"

	"github.com/theandrew168/dripfile/pkg/jsonlog"

	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

const stripePriceID = "price_1KhKCdFGWxTaRTVh9iXdyx0M"

type stripeImpl struct {
	logger     *jsonlog.Logger
	client     *client.API
	secretKey  string
	successURL string
	cancelURL  string
}

func New(logger *jsonlog.Logger, secretKey, successURL, cancelURL string) Interface {
	sc := &client.API{}
	sc.Init(secretKey, nil)

	i := stripeImpl{
		logger:     logger,
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

	i.logger.PrintInfo("stripe customer create", map[string]string{
		"customer_id": customer.ID,
	})

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

	i.logger.PrintInfo("stripe checkout_session create", map[string]string{
		"customer_id": customerID,
	})

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

	i.logger.PrintInfo("stripe subscription create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscription.Items.Data[0].ID,
	})

	return subscription.Items.Data[0].ID, nil
}

func (i *stripeImpl) CreateUsageRecord(customerID, subscriptionItemID string, quantity int64) error {
	params := stripe.UsageRecordParams{
		Quantity:         stripe.Int64(quantity),
		SubscriptionItem: stripe.String(subscriptionItemID),
	}

	_, err := i.client.UsageRecords.New(&params)
	if err != nil {
		return err
	}

	i.logger.PrintInfo("stripe usage_record create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscriptionItemID,
		"quantity":             fmt.Sprintf("%d", quantity),
	})

	return nil
}
