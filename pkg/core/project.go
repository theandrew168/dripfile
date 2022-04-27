package core

import (
	"github.com/theandrew168/dripfile/pkg/random"
)

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

func NewProjectMock() Project {
	project := NewProject(
		random.String(8),
	)
	return project
}
