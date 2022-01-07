package core

import (
	"context"
)

type Job struct {
	Pattern string
	Src     Location
	Dst     Location
	Project Project

	// readonly (from database, after creation)
	ID int
}

func NewJob(pattern string, src, dst Location, project Project) Job {
	job := Job{
		Pattern: pattern,
		Src:     src,
		Dst:     dst,
		Project: project,
	}
	return job
}

type JobStorage interface {
	CreateJob(ctx context.Context, job *Job) error
	ReadJob(ctx context.Context, id int) (Job, error)
	UpdateJob(ctx context.Context, job Job) error
	DeleteJob(ctx context.Context, job Job) error
}
