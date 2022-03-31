package database

import (
	"github.com/theandrew168/dripfile/pkg/core"
	"github.com/theandrew168/dripfile/pkg/postgres"
)

type JobStorage struct {
	db postgres.Database
}

func NewJobStorage(db postgres.Database) *JobStorage {
	s := JobStorage{
		db: db,
	}
	return &s
}

func (s *JobStorage) Create(job *core.Job) error {
	return nil
}

func (s *JobStorage) Read(id string) (core.Job, error) {
	return core.Job{}, nil
}

func (s *JobStorage) Update(job core.Job) error {
	return nil
}

func (s *JobStorage) Delete(job core.Job) error {
	return nil
}
