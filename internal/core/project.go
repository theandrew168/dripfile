package core

type Project struct {
	// readonly (from database, after creation)
	ID string

	BillingID string
}

func NewProject(billingID string) Project {
	project := Project{
		BillingID: billingID,
	}
	return project
}
