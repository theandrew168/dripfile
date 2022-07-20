package model

type Schedule struct {
	// readonly (from database, after creation)
	ID string

	Name    string
	Expr    string
	Project Project
}

// TODO: how to represent adhoc schedules (only run manually)
func NewSchedule(name, expr string, project Project) Schedule {
	schedule := Schedule{
		Name:    name,
		Expr:    expr,
		Project: project,
	}
	return schedule
}
