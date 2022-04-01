package core

type Project struct {
	// readonly (from database, after creation)
	ID string

	CustomerID         string
	SubscriptionItemID string
}

func NewProject(customerID string) Project {
	project := Project{
		CustomerID: customerID,
	}
	return project
}
