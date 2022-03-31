package billing

type PaymentGateway interface {
	CreateCustomer(email string) (string, error)
	CreateCheckoutSession(billingID string) (string, error)
}
