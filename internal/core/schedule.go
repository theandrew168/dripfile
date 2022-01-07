package core

import (
	"context"
)

// TODO: non-cron? multiple? active vs inactive?
type Schedule struct {
	Expr    string
	Project Project

	// readonly (from database, after creation)
	ID int
}

func NewSchedule(expr string, project Project) Schedule {
	schedule := Schedule{
		Expr:    expr,
		Project: project,
	}
	return schedule
}

type ScheduleStorage interface {
	CreateSchedule(ctx context.Context, schedule *Schedule) error
	ReadSchedule(ctx context.Context, id int) (Schedule, error)
	UpdateSchedule(ctx context.Context, schedule Schedule) error
	DeleteSchedule(ctx context.Context, schedule Schedule) error
}
