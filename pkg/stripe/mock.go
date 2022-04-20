package stripe

import (
	"fmt"

	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/random"
)

type mockImpl struct {
	logger *jsonlog.Logger
}

func NewMock(logger *jsonlog.Logger) Interface {
	i := mockImpl{
		logger: logger,
	}
	return &i
}

func (i *mockImpl) CreateCustomer(email string) (string, error) {
	customerID := random.String(16)

	i.logger.PrintInfo("stripe customer create", map[string]string{
		"customer_id": customerID,
	})

	return customerID, nil
}

func (i *mockImpl) CreateCheckoutSession(customerID string) (string, error) {
	sessionURL := "/billing/success?session_id=" + random.String(16)

	i.logger.PrintInfo("stripe checkout_session create", map[string]string{
		"customer_id": customerID,
	})

	return sessionURL, nil
}

func (i *mockImpl) CreateSubscription(customerID string) (string, error) {
	subscriptionItemID := random.String(16)

	i.logger.PrintInfo("stripe subscription create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscriptionItemID,
	})

	return subscriptionItemID, nil
}

func (i *mockImpl) CreateUsageRecord(customerID, subscriptionItemID string, quantity int64) error {
	i.logger.PrintInfo("stripe usage_record create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscriptionItemID,
		"quantity":             fmt.Sprintf("%d", quantity),
	})

	return nil
}
