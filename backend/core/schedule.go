package core

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
