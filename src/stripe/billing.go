package stripe

type Billing interface {
	CreateCustomer(email string) (string, error)
	CreateCheckoutSession(customerID string) (string, error)
	CreateSubscription(customerID string) (string, error)
	CreateUsageRecord(customerID, subscriptionItemID string, quantity int64) error
}
