package core

import (
	"github.com/theandrew168/dripfile/pkg/random"
)

type Schedule struct {
	// readonly (from database, after creation)
	ID string

	Name    string
	Expr    string
	Project Project
}

func NewSchedule(name, expr string, project Project) Schedule {
	schedule := Schedule{
		Name:    name,
		Expr:    expr,
		Project: project,
	}
	return schedule
}

func NewScheduleMock(project Project) Schedule {
	schedule := NewSchedule(
		random.String(8),
		random.String(8),
		project,
	)
	return schedule
}
