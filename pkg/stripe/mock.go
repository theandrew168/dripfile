package stripe

import (
	"log"

	"github.com/theandrew168/dripfile/pkg/random"
)

type mockImpl struct {
	infoLog *log.Logger
}

func NewMock(infoLog *log.Logger) Interface {
	i := mockImpl{
		infoLog: infoLog,
	}
	return &i
}

func (i *mockImpl) CreateCustomer(email string) (string, error) {
	customerID := random.String(16)

	i.infoLog.Printf("stripe.CreateCustomer:\n")
	i.infoLog.Printf("CustomerID: %v\n", customerID)

	return customerID, nil
}

func (i *mockImpl) CreateCheckoutSession(customerID string) (string, error) {
	sessionURL := "/billing/success?session_id=" + random.String(16)

	i.infoLog.Printf("stripe.CreateCheckoutSession:\n")
	i.infoLog.Printf("CustomerID: %v\n", customerID)
	i.infoLog.Printf("SessionURL: %v\n", sessionURL)

	return sessionURL, nil
}

func (i *mockImpl) CreateSubscription(customerID string) (string, error) {
	subscriptionItemID := random.String(16)

	i.infoLog.Printf("stripe.CreateSubscription:\n")
	i.infoLog.Printf("SubscriptionItemID: %v\n", subscriptionItemID)

	return subscriptionItemID, nil
}

func (i *mockImpl) CreateUsageRecord(subscriptionItemID string, quantity int64) error {
	i.infoLog.Printf("stripe.CreateUsageRecord:\n")
	i.infoLog.Printf("SubscriptionItemID: %v\n", subscriptionItemID)
	i.infoLog.Printf("Quantity: %v\n", quantity)
	return nil
}
