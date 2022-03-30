package core

type Project struct {
	// readonly (from database, after creation)
	ID string

	BillingID       string
	BillingVerified bool
}

func NewProject(billingID string) Project {
	project := Project{
		BillingID:       billingID,
		BillingVerified: false,
	}
	return project
}
