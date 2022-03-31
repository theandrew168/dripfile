package payment

import (
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/client"
)

type stripeBilling struct {
	client    *client.API
	publicKey string
	secretKey string
}

func NewStripeBilling(publicKey, secretKey string) Billing {
	sc := &client.API{}
	sc.Init(secretKey, nil)

	b := stripeBilling{
		client:    sc,
		publicKey: publicKey,
		secretKey: secretKey,
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

func (b *stripeBilling) CreateSetupIntent(billingID string) (string, error) {
	params := stripe.SetupIntentParams{
		Customer: stripe.String(billingID),
		PaymentMethodTypes: []*string{
			stripe.String("card"),
		},
	}

	intent, err := b.client.SetupIntents.New(&params)
	if err != nil {
		return "", err
	}

	return intent.ClientSecret, nil
}
