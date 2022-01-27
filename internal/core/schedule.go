package core

// TODO: non-cron? multiple? active vs inactive?
type Schedule struct {
	// readonly (from database, after creation)
	ID string

	Expr    string
	Project Project
}

func NewSchedule(expr string, project Project) Schedule {
	schedule := Schedule{
		Expr:    expr,
		Project: project,
	}
	return schedule
}

type ScheduleStorage interface {
	Create(schedule *Schedule) error
	Read(id string) (Schedule, error)
	Update(schedule Schedule) error
	Delete(schedule Schedule) error
}
