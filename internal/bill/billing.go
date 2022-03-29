package bill

type Billing interface {
	CreateCustomer(email string) (string, error)
	CreateSetupIntent(billingID string) (string, error)
}
