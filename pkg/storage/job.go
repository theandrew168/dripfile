package storage

import (
	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type Job struct {
	pg postgres.Interface
}

func NewJob(pg postgres.Interface) *Job {
	s := Job{
		pg: pg,
	}
	return &s
}

func (s *Job) Create(job *core.Job) error {
	return nil
}

func (s *Job) Read(id string) (core.Job, error) {
	return core.Job{}, nil
}

func (s *Job) Update(job core.Job) error {
	return nil
}

func (s *Job) Delete(job core.Job) error {
	return nil
}
