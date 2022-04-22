package stripe

import (
	"fmt"

	"github.com/theandrew168/dripfile/pkg/jsonlog"
	"github.com/theandrew168/dripfile/pkg/random"
)

type mockBilling struct {
	logger *jsonlog.Logger
}

func NewMockBilling(logger *jsonlog.Logger) Billing {
	b := mockBilling{
		logger: logger,
	}
	return &b
}

func (b *mockBilling) CreateCustomer(email string) (string, error) {
	customerID := random.String(16)

	b.logger.PrintInfo("stripe customer create", map[string]string{
		"customer_id": customerID,
	})

	return customerID, nil
}

func (b *mockBilling) CreateCheckoutSession(customerID string) (string, error) {
	sessionURL := "/billing/success?session_id=" + random.String(16)

	b.logger.PrintInfo("stripe checkout_session create", map[string]string{
		"customer_id": customerID,
	})

	return sessionURL, nil
}

func (b *mockBilling) CreateSubscription(customerID string) (string, error) {
	subscriptionItemID := random.String(16)

	b.logger.PrintInfo("stripe subscription create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscriptionItemID,
	})

	return subscriptionItemID, nil
}

func (b *mockBilling) CreateUsageRecord(customerID, subscriptionItemID string, quantity int64) error {
	b.logger.PrintInfo("stripe usage_record create", map[string]string{
		"customer_id":          customerID,
		"subscription_item_id": subscriptionItemID,
		"quantity":             fmt.Sprintf("%d", quantity),
	})

	return nil
}
