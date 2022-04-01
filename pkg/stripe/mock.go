package stripe

import (
	"github.com/theandrew168/dripfile/pkg/log"
	"github.com/theandrew168/dripfile/pkg/random"
)

type mockImpl struct {
	logger log.Logger
}

func NewMock(logger log.Logger) Interface {
	i := mockImpl{
		logger: logger,
	}
	return &i
}

func (i *mockImpl) CreateCustomer(email string) (string, error) {
	customerID := random.String(16)

	i.logger.Info("stripe.CreateCustomer:\n")
	i.logger.Info("CreateCustomer:\n")
	i.logger.Info("CustomerID: %s\n", customerID)

	return customerID, nil
}

func (i *mockImpl) CreateCheckoutSession(customerID string) (string, error) {
	sessionURL := "/billing/success?session_id=" + random.String(16)

	i.logger.Info("stripe.CreateCheckoutSession:\n")
	i.logger.Info("CustomerID: %s\n", customerID)
	i.logger.Info("SessionURL: %s\n", sessionURL)

	return sessionURL, nil
}

func (i *mockImpl) CreateSubscription(customerID string) (string, error) {
	i.logger.Info("stripe.CreateSubscription:\n")

	return "todo", nil
}
