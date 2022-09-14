package model

type Schedule struct {
	// readonly (from database, after creation)
	ID string

	Name string
	Expr string
}

// TODO: how to represent adhoc schedules (only run manually)
func NewSchedule(name, expr string) Schedule {
	schedule := Schedule{
		Name: name,
		Expr: expr,
	}
	return schedule
}
