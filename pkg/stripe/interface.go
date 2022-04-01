package stripe

type Interface interface {
	CreateCustomer(email string) (string, error)
	CreateCheckoutSession(customerID string) (string, error)
	CreateSubscription(customerID string) (string, error)
}
