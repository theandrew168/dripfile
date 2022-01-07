package core

import (
	"context"
)

// TODO: non-cron? multiple? active vs inactive?
type Schedule struct {
	Name string
	Expr string

	// readonly (from database, after creation)
	ID int
}

func NewSchedule(name, expr string) Schedule {
	schedule := Schedule{
		Name: name,
		Expr: expr,
	}
	return schedule
}

type ScheduleStorage interface {
	CreateSchedule(ctx context.Context, schedule *Schedule) error
	ReadSchedule(ctx context.Context, id int) (Schedule, error)
	UpdateSchedule(ctx context.Context, schedule Schedule) error
	DeleteSchedule(ctx context.Context, schedule Schedule) error
}
