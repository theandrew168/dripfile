package bill

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

type stripeBilling struct {
	client       *client.API
	stripeAPIKey string
}

func NewStripeBilling(stripeAPIKey string) Billing {
	sc := &client.API{}
	sc.Init(stripeAPIKey, nil)

	b := stripeBilling{
		client:       sc,
		stripeAPIKey: stripeAPIKey,
	}
	return &b
}

func (b *stripeBilling) CreateCustomer(email string) (string, error) {
	params := stripe.CustomerParams{
		Email: stripe.String(email),
	}
	customer, err := b.client.Customers.New(&params)
	if err != nil {
		return "", err
	}
	return customer.ID, nil
}
