package core

import (
	"context"
)

type Job struct {
	Src Location
	Dst Location

	// readonly (from database, after creation)
	ID int
}

func NewJob(src, dst Location) Job {
	job := Job{
		Src: src,
		Dst: dst,
	}
	return job
}

type JobStorage interface {
	CreateJob(ctx context.Context, job *Job) error
	ReadJob(ctx context.Context, id int) (Job, error)
	UpdateJob(ctx context.Context, job Job) error
	DeleteJob(ctx context.Context, job Job) error
}
